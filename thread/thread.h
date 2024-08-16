// thread.h
#ifndef THREAD_H
#define THREAD_H

// set_thread_affinity sets the affinity of the current thread to the given cpuid.
int set_thread_affinity(int cpuid);

// thread_func_t is a function pointer type that takes a void pointer as an argument and returns void.
typedef void (*thread_func_t)(void *);

// thread_data is a struct that contains the thread index, a thread function, and an argument.
void create_thread(unsigned thread_index, thread_func_t f, void *arg);

// hardware_concurrency returns the number of hardware threads available on the system.
unsigned hardware_concurrency();

#endif // THREAD_H