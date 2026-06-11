package me.killkiss.honokactrl.service

import android.app.Notification
import android.app.NotificationChannel
import android.app.NotificationManager
import android.app.PendingIntent
import android.app.Service
import android.content.Context
import android.content.Intent
import android.os.Build
import android.os.IBinder
import android.os.PowerManager
import android.util.Log
import androidx.core.app.NotificationCompat
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.SupervisorJob
import kotlinx.coroutines.cancel
import kotlinx.coroutines.launch
import me.killkiss.honokactrl.MainActivity
import me.killkiss.honokactrl.R
import me.killkiss.honokactrl.jni.NativeBridge

class HonokaServerService : Service() {
    companion object {
        private const val TAG = "HonokaServer"
        private var wakeLock: PowerManager.WakeLock? = null
        private const val CHANNEL_ID = "honoka_server"
        private const val NOTIFICATION_ID = 1001

        const val ACTION_START = "me.killkiss.honokactrl.action.START"
        const val ACTION_STOP = "me.killkiss.honokactrl.action.STOP"
        const val EXTRA_RUNTIME_ROOT = "runtime_root"

        fun startIntent(context: Context, runtimeRoot: String): Intent {
            return Intent(context, HonokaServerService::class.java).apply {
                action = ACTION_START
                putExtra(EXTRA_RUNTIME_ROOT, runtimeRoot)
            }
        }

        fun stopIntent(context: Context): Intent {
            return Intent(context, HonokaServerService::class.java).apply {
                action = ACTION_STOP
            }
        }
    }

    private val scope = CoroutineScope(SupervisorJob() + Dispatchers.IO)
    private var runtimeRoot: String? = null

    override fun onCreate() {
        super.onCreate()
        Log.i(TAG, "service created")
        ensureChannel()
    }

    override fun onStartCommand(intent: Intent?, flags: Int, startId: Int): Int {
        Log.i(TAG, "onStartCommand: action=${intent?.action}, startId=$startId")
        when (intent?.action) {
            ACTION_START -> {
                startForeground(NOTIFICATION_ID, buildNotification("正在启动服务"))
                acquireWakeLock()
                runtimeRoot = intent.getStringExtra(EXTRA_RUNTIME_ROOT)
                if (runtimeRoot.isNullOrBlank()) {
                    Log.e(TAG, "missing runtime root")
                    updateNotification("缺少运行目录")
                    stopSelf()
                    return START_REDELIVER_INTENT
                }

                scope.launch {
                    Log.i(TAG, "native start begin, runtimeRoot=$runtimeRoot")
                    val result = NativeBridge.start(runtimeRoot!!)
                    if (result.isSuccess) {
                        Log.i(TAG, "native start success")
                        updateNotification("服务运行中")
                    } else {
                        Log.e(TAG, "native start failed", result.exceptionOrNull())
                        updateNotification(result.exceptionOrNull()?.message ?: "启动失败")
                        releaseWakeLock()
                        stopSelf()
                    }
                }
            }

            ACTION_STOP -> {
                scope.launch {
                    Log.i(TAG, "native stop begin")
                    NativeBridge.stop()
                    Log.i(TAG, "native stop finished")
                    updateNotification("服务已停止")
                    stopForeground(STOP_FOREGROUND_REMOVE)
                    releaseWakeLock()
                    stopSelf()
                }
                return START_NOT_STICKY
            }
        }
        return START_REDELIVER_INTENT
    }

    override fun onDestroy() {
        Log.i(TAG, "service destroyed")
        scope.cancel()
        releaseWakeLock()
        super.onDestroy()
    }

    override fun onBind(intent: Intent?): IBinder? = null

    private fun buildNotification(content: String): Notification {
        val openIntent = PendingIntent.getActivity(
            this,
            0,
            Intent(this, MainActivity::class.java).apply {
                flags = Intent.FLAG_ACTIVITY_SINGLE_TOP or Intent.FLAG_ACTIVITY_CLEAR_TOP
            },
            PendingIntent.FLAG_UPDATE_CURRENT or PendingIntent.FLAG_IMMUTABLE,
        )

        return NotificationCompat.Builder(this, CHANNEL_ID)
            .setSmallIcon(R.drawable.ic_notification_server)
            .setContentTitle(getString(R.string.app_name))
            .setContentText(content)
            .setContentIntent(openIntent)
            .setOngoing(true)
            .build()
    }

    private fun updateNotification(content: String) {
        val manager = getSystemService(Context.NOTIFICATION_SERVICE) as NotificationManager
        manager.notify(NOTIFICATION_ID, buildNotification(content))
    }

    private fun ensureChannel() {
        if (Build.VERSION.SDK_INT < Build.VERSION_CODES.O) {
            return
        }
        val manager = getSystemService(Context.NOTIFICATION_SERVICE) as NotificationManager
        val channel = NotificationChannel(
            CHANNEL_ID,
            "Honoka Server",
            NotificationManager.IMPORTANCE_LOW,
        ).apply {
            description = "honoka-chan background service"
        }
        manager.createNotificationChannel(channel)
    }

    private fun acquireWakeLock() {
        if (wakeLock?.isHeld == true) {
            return
        }
        val powerManager = getSystemService(Context.POWER_SERVICE) as PowerManager
        wakeLock = powerManager.newWakeLock(PowerManager.PARTIAL_WAKE_LOCK, "honoka:server").apply {
            setReferenceCounted(false)
            acquire()
        }
        Log.i(TAG, "wakelock acquired")
    }

    private fun releaseWakeLock() {
        wakeLock?.takeIf { it.isHeld }?.release()
        wakeLock = null
        Log.i(TAG, "wakelock released")
    }
}
