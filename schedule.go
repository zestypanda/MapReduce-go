package mapreduce

import "fmt"
import "sync"
//import "log"

//
// schedule() starts and waits for all tasks in the given phase (Map
// or Reduce). the mapFiles argument holds the names of the files that
// are the inputs to the map phase, one per map task. nReduce is the
// number of reduce tasks. the registerChan argument yields a stream
// of registered workers; each item is the worker's RPC address,
// suitable for passing to call(). registerChan will yield all
// existing registered workers (if any) and new ones as they register.
//
func schedule(jobName string, mapFiles []string, nReduce int, phase jobPhase, registerChan chan string) {
	var ntasks int
	var n_other int // number of inputs (for reduce) or outputs (for map)
	switch phase {
	case mapPhase:
		ntasks = len(mapFiles)
		n_other = nReduce
	case reducePhase:
		ntasks = nReduce
		n_other = len(mapFiles)
	}

	fmt.Printf("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, n_other)

	var wait_group sync.WaitGroup;

    for i := 0; i < ntasks; i++ {
        wait_group.Add(1)
        // Every task has its own DoTaskArgs
        var taskArgs DoTaskArgs
        taskArgs.JobName = jobName
        taskArgs.Phase = phase
        taskArgs.TaskNumber = i
        taskArgs.NumOtherPhase = n_other
        if (phase == mapPhase) {
            taskArgs.File = mapFiles[i]
        }
        go func() {
            // fmt.Printf("Processing: %dth task\n", i)
            defer wait_group.Done()
            for {
                wk := <-registerChan
                if (call(wk, "Worker.DoTask", &taskArgs, nil) == true) {
                    // put back worker when finished
                    go func() {registerChan <- wk}()
                    break
                }
            }           
        }()
    }
    wait_group.Wait()
    
	// All ntasks tasks have to be scheduled on workers, and only once all of
	// them have been completed successfully should the function return.
	// Remember that workers may fail, and that any given worker may finish
	// multiple tasks.
	//
	// TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO
	//
	fmt.Printf("Schedule: %v phase done\n", phase)
}
