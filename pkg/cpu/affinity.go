package cpu

/*
#ifdef _WIN32
# include <windows.h>
#elif defined(__linux__)
# ifndef _GNU_SOURCE
#  define _GNU_SOURCE
# endif
# include <sched.h>
# include <unistd.h>
# include <pthread.h>

#else

#endif

#include <stdio.h>
#include <stdlib.h>

int _set_process_affinity_(int start, int step, int count, int pid)
{
    if (start < 0 || step < 1 || count < 1) {
        return 1;
    }

#ifdef _WIN32

    return 0;
    // https://learn.microsoft.com/en-us/windows/win32/procthread/cpu-sets
    // const int max_count = 4096;
    // if (count > max_count) {
    //     return 1;
    // }
    // ULONG cpuSetIds[max_count];
    // for (int index = 0; index < count; index++) {
    //     int cpu = start + index * step;
    //     cpuSetIds[index] = (ULONG) cpu;
    // }
    // if (!SetProcessDefaultCpuSets(GetCurrentProcess(), cpuSetIds, (ULONG) count)) {
    //     auto code = ::GetLastError();
    //     printf("SetProcessDefaultCpuSets=%d\n", code);
    //     return -1;
    // }
    // return 0;
#elif defined(__linux__)
    const int nrcpus = sysconf(_SC_NPROCESSORS_CONF);
    cpu_set_t* pset = CPU_ALLOC(nrcpus);//, [](cpu_set_t *set) { if (set) CPU_FREE(set); });
    if (!pset) {
        perror("CPU_ALLOC");
        return -1;
    }
    size_t cpusize = CPU_ALLOC_SIZE(nrcpus);
    CPU_ZERO_S(cpusize, pset);
    for (int index = 0; index < count; index++) {
        int cpu = start + index * step;
        CPU_SET_S(cpu, cpusize, pset);
    }
	if(pid == 0){
		pid = getpid();
	}
    if (sched_setaffinity(pid, cpusize, pset) < 0) {
        perror("sched_setaffinity");
		CPU_FREE(pset);
        return -1;
    }
	CPU_FREE(pset);
    return 0;
#else
    // Not Support
    return 1;
#endif
}
*/
import "C"

func SetProcessAffinity(start, step, count int, pid int) {
	C._set_process_affinity_(C.int(start), C.int(step), C.int(count), C.int(pid))
}
