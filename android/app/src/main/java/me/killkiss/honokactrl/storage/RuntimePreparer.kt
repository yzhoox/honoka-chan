package me.killkiss.honokactrl.storage

import android.content.Context
import android.os.Build
import android.system.Os
import android.util.Log
import java.io.File
import java.io.FileOutputStream
import java.io.IOException
import java.nio.file.Files
import java.security.MessageDigest
import java.util.zip.ZipInputStream

class RuntimePreparer(private val context: Context) {
    fun prepare(dataDirPath: String, forceRedeploy: Boolean = false): Result<RuntimeInfo> {
        return runCatching {
            val runtimeRoot = File(context.filesDir, RUNTIME_DIR_NAME)
            val bundleHash = calculateBundleHash()
            Log.i(TAG, "prepare runtime root=${runtimeRoot.absolutePath}, dataDir=$dataDirPath, forceRedeploy=$forceRedeploy")
            if (forceRedeploy || shouldRedeployBundle(runtimeRoot, bundleHash)) {
                redeployBundle(runtimeRoot, bundleHash)
            }
            mountAndroidDir(runtimeRoot, File(dataDirPath.ifBlank { SettingsStore.DEFAULT_DATA_DIR_PATH }))
            RuntimeInfo(
                runtimeRoot = runtimeRoot,
                mountedDataDir = File(dataDirPath.ifBlank { SettingsStore.DEFAULT_DATA_DIR_PATH }),
                configFile = File(runtimeRoot, "config.json"),
                deployedBundleHash = bundleHash,
            )
        }
    }

    private fun redeployBundle(runtimeRoot: File, bundleHash: String) {
        Log.i(TAG, "redeploy bundle into ${runtimeRoot.absolutePath}")
        val preservedFiles = PRESERVED_RUNTIME_FILES.mapNotNull { relativePath ->
            File(runtimeRoot, relativePath).takeIf { it.exists() }?.let { relativePath to it.readBytes() }
        }.toMap()

        if (runtimeRoot.exists()) {
            runtimeRoot.deleteRecursively()
        }
        runtimeRoot.mkdirs()

        context.assets.open(BUNDLE_ASSET_NAME).use { input ->
            ZipInputStream(input).use { zip ->
                var entry = zip.nextEntry
                while (entry != null) {
                    val outFile = File(runtimeRoot, entry.name)
                    if (preservedFiles.containsKey(entry.name)) {
                        zip.closeEntry()
                        entry = zip.nextEntry
                        continue
                    }
                    if (entry.isDirectory) {
                        outFile.mkdirs()
                    } else {
                        outFile.parentFile?.mkdirs()
                        FileOutputStream(outFile).use { output ->
                            zip.copyTo(output)
                        }
                    }
                    zip.closeEntry()
                    entry = zip.nextEntry
                }
            }
        }

        preservedFiles.forEach { (relativePath, data) ->
            File(runtimeRoot, relativePath).apply {
                parentFile?.mkdirs()
                writeBytes(data)
            }
        }

        File(runtimeRoot, BUNDLE_MARKER_FILE).writeText(bundleHash)
        Log.i(TAG, "bundle deployed successfully")
    }

    private fun shouldRedeployBundle(runtimeRoot: File, bundleHash: String): Boolean {
        val markerFile = File(runtimeRoot, BUNDLE_MARKER_FILE)
        if (!markerFile.exists()) {
            return true
        }

        val deployedHash = runCatching { markerFile.readText().trim() }.getOrDefault("")
        val changed = deployedHash != bundleHash
        if (changed) {
            Log.i(TAG, "runtime bundle hash changed, redeploy required")
        }
        return changed
    }

    private fun calculateBundleHash(): String {
        val digest = MessageDigest.getInstance("SHA-256")
        context.assets.open(BUNDLE_ASSET_NAME).use { input ->
            val buffer = ByteArray(DEFAULT_BUFFER_SIZE)
            while (true) {
                val read = input.read(buffer)
                if (read <= 0) {
                    break
                }
                digest.update(buffer, 0, read)
            }
        }
        return digest.digest().joinToString("") { "%02x".format(it) }
    }

    private fun mountAndroidDir(runtimeRoot: File, targetDir: File) {
        if (!targetDir.exists()) {
            targetDir.mkdirs()
        }
        Log.i(TAG, "mount Android dir target=${targetDir.absolutePath}")

        val androidDir = File(runtimeRoot, "static/Android")
        androidDir.parentFile?.mkdirs()

        if (isSymlink(androidDir)) {
            androidDir.delete()
        } else if (androidDir.exists()) {
            androidDir.deleteRecursively()
        }

        Os.symlink(targetDir.absolutePath, androidDir.absolutePath)
        Log.i(TAG, "symlink created: ${androidDir.absolutePath} -> ${targetDir.absolutePath}")
    }

    private fun isSymlink(file: File): Boolean {
        return if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            Files.isSymbolicLink(file.toPath())
        } else {
            try {
                file.canonicalFile != file.absoluteFile
            } catch (_: IOException) {
                false
            }
        }
    }

    companion object {
        private const val TAG = "HonokaServer"
        const val BUNDLE_ASSET_NAME = "honoka_runtime.zip"
        const val RUNTIME_DIR_NAME = "honoka-runtime"
        private const val BUNDLE_MARKER_FILE = ".bundle_ready"
        private val PRESERVED_RUNTIME_FILES = listOf(
            "config.json",
            "assets/data.db",
            "assets/data.db-wal",
            "assets/data.db-shm",
        )

        fun readDeployedBundleHash(runtimeRoot: File): String {
            return runCatching {
                File(runtimeRoot, BUNDLE_MARKER_FILE).takeIf { it.exists() }?.readText()?.trim().orEmpty()
            }.getOrDefault("")
        }
    }
}

data class RuntimeInfo(
    val runtimeRoot: File,
    val mountedDataDir: File,
    val configFile: File,
    val deployedBundleHash: String,
)
