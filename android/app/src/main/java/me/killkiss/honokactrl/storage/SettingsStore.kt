package me.killkiss.honokactrl.storage

import android.content.Context

class SettingsStore(context: Context) {
    private val prefs = context.getSharedPreferences(PREFS_NAME, Context.MODE_PRIVATE)

    fun loadDataDirPath(): String = prefs.getString(KEY_DATA_DIR_PATH, DEFAULT_DATA_DIR_PATH) ?: DEFAULT_DATA_DIR_PATH

    fun saveDataDirPath(path: String) {
        prefs.edit().putString(KEY_DATA_DIR_PATH, path).apply()
    }

    companion object {
        private const val PREFS_NAME = "honoka_control"
        private const val KEY_DATA_DIR_PATH = "data_dir_path"
        const val DEFAULT_DATA_DIR_PATH = "/sdcard/Download/data"
    }
}
