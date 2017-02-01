package main

import "etec3702"
import "math/rand"
import "runtime"

var t_on_table int
var p_on_table int
var m_on_table int
var s_pm_sem = etec3702.NewSemaphore(0)
var s_pt_sem = etec3702.NewSemaphore(0)
var s_mt_sem = etec3702.NewSemaphore(0)
var table_t = etec3702.NewSemaphore(0)
var table_p = etec3702.NewSemaphore(0)
var table_m = etec3702.NewSemaphore(0)
var pusherlock = etec3702.NewSemaphore(1)

func matches_smoker(){
    for {
        etec3702.Delay()
        s_pt_sem.Acquire()
        etec3702.Output("m smoke")
        etec3702.Delay()
    }
}

func tobacco_smoker(){
    for {
        etec3702.Delay()
        s_pm_sem.Acquire()
        etec3702.Output("t smoke")
        etec3702.Delay()
    }
}

func paper_smoker(){
    for {
        etec3702.Delay()
        s_mt_sem.Acquire()
        etec3702.Output("p smoke")
        etec3702.Delay()
    }
}

func tobacco_pusher(){
    for{
        table_t.Acquire()
        pusherlock.Acquire()
        if p_on_table > 0 {
            p_on_table--
            s_pt_sem.Release()
        } else if m_on_table > 0 {
            m_on_table--
            s_mt_sem.Release()
        } else{
            t_on_table++
        }
        pusherlock.Release()
    }
}


func paper_pusher(){
    for{
        table_p.Acquire()
        pusherlock.Acquire()
        if t_on_table > 0 {
            t_on_table--
            s_pt_sem.Release()
        } else if m_on_table > 0 {
            m_on_table--
            s_pm_sem.Release()
        }  else{
            p_on_table++
        }
        pusherlock.Release()
    }
}


func match_pusher(){
    for{
        table_m.Acquire()
        pusherlock.Acquire()
        if t_on_table>0 {
            t_on_table--
            s_mt_sem.Release()
        } else if p_on_table>0 {
            p_on_table--
            s_pm_sem.Release()
        } else{
            m_on_table++
        }
        pusherlock.Release()
    }
}
        
func agent(){
    for {
        switch(rand.Intn(3)){
            case 0:
                etec3702.Output("put matches on table")
                table_m.Release()
                etec3702.Output("put paper on table")
                table_p.Release()

            case 1:
                etec3702.Output("put tobacco on table")
                table_t.Release()
                etec3702.Output("put paper on table")
                table_p.Release()
            case 2:
                etec3702.Output("put paper on table")
                table_t.Release()
                etec3702.Output("put matches on table")
                table_m.Release()
        }
        etec3702.Delay()
    }
}

        
func main(){
    go agent()
    go match_pusher()
    go paper_pusher()
    go tobacco_pusher()
    go matches_smoker()
    go paper_smoker()
    go tobacco_smoker()
    
    runtime.Goexit()
}

    
