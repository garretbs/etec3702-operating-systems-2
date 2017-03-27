
import queue
import threading

queues={}


def actor1(Q):
    while 1:
        m = Q.get()
        if m == "quit":
            return
        print("A1: I got this:",m)
        queues["A2"].put(m+"!")

def actor2(Q):
    while 1:
        m = Q.get()
        if m == "quit":
            return
        print("A2: I got this:",m)
        queues["A3"].put(m+"?")

def actor3(Q):
    while 1:
        m = Q.get()
        print("A3: I got this:",m)
        if len(m) > 10:
            print("Stopping")
            queues["A1"].put("quit")
            queues["A2"].put("quit")
            return
        queues["A1"].put(m+".")


threads=[]
for name,func in [ ("A1",actor1),("A2",actor2),("A3",actor3)]:
    queues[name]=queue.Queue()
    t1 = threading.Thread(target=func, args=(queues[name],))
    threads.append(t1)
    
for t in threads:
    t.start()

queues["A1"].put("GO")

for t in threads:
    t.join()
    
