#include <jni.h>
#include <dlfcn.h>
#include <mutex>
#include <string>

namespace {
using server_start_fn = char* (*)(const char*);
using server_stop_fn = char* (*)();
using server_status_fn = char* (*)();
using server_health_fn = char* (*)();
using server_reload_fn = char* (*)();
using server_free_string_fn = void (*)(char*);

std::once_flag g_init_once;
void* g_handle = nullptr;
server_start_fn g_start = nullptr;
server_stop_fn g_stop = nullptr;
server_status_fn g_status = nullptr;
server_health_fn g_health = nullptr;
server_reload_fn g_reload = nullptr;
server_free_string_fn g_free_string = nullptr;
std::string g_error;

void InitServerSymbols() {
    g_handle = dlopen("libhonokachan.so", RTLD_NOW | RTLD_GLOBAL);
    if (g_handle == nullptr) {
        g_error = dlerror();
        return;
    }

    g_start = reinterpret_cast<server_start_fn>(dlsym(g_handle, "ServerStart"));
    g_stop = reinterpret_cast<server_stop_fn>(dlsym(g_handle, "ServerStop"));
    g_status = reinterpret_cast<server_status_fn>(dlsym(g_handle, "ServerStatusJSON"));
    g_health = reinterpret_cast<server_health_fn>(dlsym(g_handle, "ServerHealthJSON"));
    g_reload = reinterpret_cast<server_reload_fn>(dlsym(g_handle, "ServerReload"));
    g_free_string = reinterpret_cast<server_free_string_fn>(dlsym(g_handle, "ServerFreeString"));

    if (g_start == nullptr || g_stop == nullptr || g_status == nullptr || g_health == nullptr ||
        g_reload == nullptr || g_free_string == nullptr) {
        g_error = "missing exported symbols from libhonokachan.so";
    }
}

const char* EnsureServerReady() {
    std::call_once(g_init_once, InitServerSymbols);
    if (!g_error.empty()) {
        return g_error.c_str();
    }
    return nullptr;
}

jstring CallStringResult(JNIEnv* env, char* result) {
    if (result == nullptr) {
        return nullptr;
    }

    jstring java_string = env->NewStringUTF(result);
    g_free_string(result);
    return java_string;
}
}  // namespace

extern "C"
JNIEXPORT jstring JNICALL
Java_me_killkiss_honokactrl_jni_NativeBridge_nativeStart(JNIEnv* env, jobject, jstring work_dir) {
    const char* error = EnsureServerReady();
    if (error != nullptr) {
        return env->NewStringUTF(error);
    }

    const char* work_dir_chars = env->GetStringUTFChars(work_dir, nullptr);
    char* result = g_start(work_dir_chars);
    env->ReleaseStringUTFChars(work_dir, work_dir_chars);
    return CallStringResult(env, result);
}

extern "C"
JNIEXPORT jstring JNICALL
Java_me_killkiss_honokactrl_jni_NativeBridge_nativeStop(JNIEnv* env, jobject) {
    const char* error = EnsureServerReady();
    if (error != nullptr) {
        return env->NewStringUTF(error);
    }

    return CallStringResult(env, g_stop());
}

extern "C"
JNIEXPORT jstring JNICALL
Java_me_killkiss_honokactrl_jni_NativeBridge_nativeStatusJson(JNIEnv* env, jobject) {
    const char* error = EnsureServerReady();
    if (error != nullptr) {
        std::string json = std::string("{\"running\":false,\"last_error\":\"") + error + "\"}";
        return env->NewStringUTF(json.c_str());
    }

    char* result = g_status();
    if (result == nullptr) {
        return env->NewStringUTF("{\"running\":false,\"last_error\":\"status unavailable\"}");
    }
    jstring java_string = env->NewStringUTF(result);
    g_free_string(result);
    return java_string;
}

extern "C"
JNIEXPORT jstring JNICALL
Java_me_killkiss_honokactrl_jni_NativeBridge_nativeHealthJson(JNIEnv* env, jobject) {
    const char* error = EnsureServerReady();
    if (error != nullptr) {
        std::string json = std::string("{\\\"status\\\":\\\"error\\\",\\\"message\\\":\\\"") + error + "\\\"}";
        return env->NewStringUTF(json.c_str());
    }

    char* result = g_health();
    if (result == nullptr) {
        return env->NewStringUTF("{\\\"status\\\":\\\"error\\\",\\\"message\\\":\\\"health unavailable\\\"}");
    }
    jstring java_string = env->NewStringUTF(result);
    g_free_string(result);
    return java_string;
}

extern "C"
JNIEXPORT jstring JNICALL
Java_me_killkiss_honokactrl_jni_NativeBridge_nativeReload(JNIEnv* env, jobject) {
    const char* error = EnsureServerReady();
    if (error != nullptr) {
        return env->NewStringUTF(error);
    }

    return CallStringResult(env, g_reload());
}
