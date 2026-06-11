package me.killkiss.honokactrl.jni

import org.json.JSONObject

data class NativeServerStatus(
    val running: Boolean = false,
    val workDir: String = "",
    val listenPort: String = "",
    val localUrl: String = "",
    val startedAt: String = "",
    val lastError: String = "",
)

data class NativeReloadResult(
    val success: Boolean,
    val message: String,
    val reloadedAt: String = "",
)

data class NativeHealthInfo(
    val status: String = "",
    val appName: String = "",
    val version: String = "",
    val startedAt: String = "",
    val uptimeSeconds: Long = 0,
    val lastReloadAt: String = "",
    val listenPort: String = "",
    val cdnServer: String = "",
    val reloadTokenConfigured: Boolean = false,
    val mainDb: String = "",
    val userDb: String = "",
    val message: String = "",
)

object NativeBridge {
    private val loadError: String? by lazy {
        runCatching { System.loadLibrary("honoka_android") }
            .exceptionOrNull()
            ?.message
    }

    private external fun nativeStart(workDir: String): String?
    private external fun nativeStop(): String?
    private external fun nativeStatusJson(): String
    private external fun nativeHealthJson(): String
    private external fun nativeReload(): String?

    fun start(workDir: String): Result<Unit> {
        val error = ensureLoaded() ?: nativeStart(workDir)
        return if (error == null) Result.success(Unit) else Result.failure(IllegalStateException(error))
    }

    fun stop(): Result<Unit> {
        val error = ensureLoaded() ?: nativeStop()
        return if (error == null) Result.success(Unit) else Result.failure(IllegalStateException(error))
    }

    fun status(): NativeServerStatus {
        val error = ensureLoaded()
        if (error != null) {
            return NativeServerStatus(lastError = error)
        }

        return runCatching {
            val json = JSONObject(nativeStatusJson())
            NativeServerStatus(
                running = json.optBoolean("running"),
                workDir = json.optString("work_dir"),
                listenPort = json.optString("listen_port"),
                localUrl = json.optString("local_url"),
                startedAt = json.optString("started_at"),
                lastError = json.optString("last_error"),
            )
        }.getOrElse {
            NativeServerStatus(lastError = it.message ?: "failed to parse native status")
        }
    }

    private fun healthJson(): String {
        val error = ensureLoaded()
        if (error != null) {
            return """{"status":"error","message":${JSONObject.quote(error)}}"""
        }
        return nativeHealthJson()
    }

    fun health(): NativeHealthInfo {
        val raw = healthJson()
        return runCatching {
            val json = JSONObject(raw)
            NativeHealthInfo(
                status = json.optString("status"),
                appName = json.optString("app_name"),
                version = json.optString("version"),
                startedAt = json.optString("started_at"),
                uptimeSeconds = json.optLong("uptime_seconds"),
                lastReloadAt = json.optString("last_reload_at"),
                listenPort = json.optString("listen_port"),
                cdnServer = json.optString("cdn_server"),
                reloadTokenConfigured = json.optBoolean("reload_token_configured"),
                mainDb = json.optString("main_db"),
                userDb = json.optString("user_db"),
                message = json.optString("message"),
            )
        }.getOrElse {
            NativeHealthInfo(
                status = "error",
                message = it.message ?: "failed to parse health response",
            )
        }
    }

    fun reload(): NativeReloadResult {
        val error = ensureLoaded()
        if (error != null) {
            return NativeReloadResult(success = false, message = error)
        }

        val result = nativeReload()
        if (result == null) {
            return NativeReloadResult(success = true, message = "configuration reloaded")
        }

        return runCatching {
            val json = JSONObject(result)
            NativeReloadResult(
                success = json.optString("status") == "ok",
                message = json.optString("message").ifBlank { "reload failed" },
                reloadedAt = json.optString("reloaded_at"),
            )
        }.getOrElse {
            NativeReloadResult(success = false, message = it.message ?: "reload failed")
        }
    }

    private fun ensureLoaded(): String? = loadError
}
