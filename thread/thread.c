// thread.c
#define _GNU_SOURCE

#include "thread.h"
#include <stdio.h>
#include <stdlib.h>
#if defined(_WIN32) || defined(__CYGWIN__)
  #include <Windows.h>
#else
  #ifdef __APPLE__
    #include <mach/thread_act.h>
    #include <mach/thread_policy.h>
  #endif
  #include <pthread.h>
  #include <unistd.h>
  #include <sys/types.h>
  #include <sys/sysctl.h>

  #include <sched.h> // Ensure this header is included for CPU_ZERO and CPU_SET
#endif


int
set_thread_affinity(int cpuid)
{
#if defined(_WIN32) || defined(__CYGWIN__)
    int numThread = GetMaximumProcessorCount(0);
    int core = cpuid % numThread;
    int node = cpuid / numThread;
    GROUP_AFFINITY affinity = { 0 };
    affinity.Mask = (KAFFINITY)(1ULL << core);
    affinity.Group = node;
    BOOL success = SetThreadGroupAffinity(GetCurrentThread(), &affinity, NULL);
    if (!success) {
        //std::cerr << "Failed to set thread affinity for thread id: " << cpuid << " (error=" << GetLastError() << ")" << std::endl;
        return GetLastError();
    }
    return 0;
    //     std::cerr << "Failed to set thread affinity for thread id: " << cpuId << " (error=" << GetLastError() << ")" << std::endl;
    //     return 1;
    // }
    // return 0;

    //thread = reinterpret_cast<std::thread::native_handle_type>(GetCurrentThread());
#elif defined(__APPLE__)
    thread_port_t mach_thread;
    thread_affinity_policy_data_t policy = { cpuid };
    mach_thread = pthread_mach_thread_np(pthread_self());
    return thread_policy_set(mach_thread, THREAD_AFFINITY_POLICY,
            (thread_policy_t)&policy, 1);
    //thread = static_cast<std::thread::native_handle_type>(pthread_self());
#elif !defined(__OpenBSD__) && !defined(__FreeBSD__) && !defined(__ANDROID__) && !defined(__NetBSD__)
    cpu_set_t cs;
    CPU_ZERO(&cs);
    CPU_SET(cpuid, &cs);

    int rc = pthread_setaffinity_np(pthread_self(), sizeof(cpu_set_t), &cs);
    if (rc != 0) {
        // std::cerr << "Failed to set thread affinity for thread id: " << cpuid << " (error=" << rc << ")" << std::endl;
        return rc;
    }
    return 0;
#endif
    return -1;
}

typedef struct {
    unsigned thread_index;
    thread_func_t f;
    void *arg;
} thread_data;

void thread_func_wrapper(void *arg);
void threadFunc(void *arg);

void thread_run(void *arg) {
    thread_data *data = (thread_data *)arg;
    if (data->thread_index != -1){
        set_thread_affinity(data->thread_index);
    }
    data->f(data->arg);
    free(data);
}

void create_thread(unsigned thread_index, thread_func_t f, void *arg) {
    thread_data *data = (thread_data *)malloc(sizeof(thread_data));
    data->thread_index = thread_index;
    data->f = f;
    data->arg = arg;

    // call arg
//    threadFunc(arg);

    // Platform-specific thread creation code
    #if defined(_WIN32) || defined(__CYGWIN__)
        HANDLE thread = CreateThread(NULL, 0, (LPTHREAD_START_ROUTINE)thread_run, data, 0, NULL);
        if (thread != NULL) {
            CloseHandle(thread);
        }
    #else
        pthread_t thread;
        pthread_create(&thread, NULL, (void *(*)(void *))thread_run, data);
        pthread_detach(thread);
    #endif
}

unsigned hardware_concurrency() {
#if defined(_WIN32) || defined(__CYGWIN__)
    SYSTEM_INFO sysinfo;
    GetSystemInfo(&sysinfo);
    return sysinfo.dwNumberOfProcessors;
#else
    return sysconf(_SC_NPROCESSORS_ONLN);
#endif
}