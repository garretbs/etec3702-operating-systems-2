import random
import threading
import queue

wins_queue = []
plays_queue = []
total_plays = 10000
num_sub_threads = 4

def create_sub_threads(strategy):
	wins = 0
	plays = 0
	
	sub_threads = []
	for i in range(num_sub_threads):
		sub_thread = threading.Thread(target = monty_hall, args=[strategy])
		sub_thread.start()
		sub_threads.append(sub_thread)
	for sub_thread in sub_threads:
		sub_thread.join()
		
	while plays < total_plays:
		plays += plays_queue[strategy].get()
	while wins_queue[strategy].qsize() > 0:
		wins += wins_queue[strategy].get()
		
	print("Strategy",strategy,wins/total_plays*100,"%")


def monty_hall(strategy):
	sub_plays = (int) (total_plays/num_sub_threads)
	for trials in range(sub_plays):
		#the winning door
		winning = random.randrange(3)
		#the player's choice
		pick = random.randrange(3)
		#the door Monty opens
		open_ = random.choice( list(set([0,1,2]) - set([winning]) - set([pick]) ) )
		
		#the remaining door
		remaining = (set([0,1,2]) - set([pick]) - set([open_])).pop()
		
		if strategy == 0:
			#strategy 0: Stand pat.
			pass 
		elif strategy == 1:
			#strategy 1: Flip a coin, choose between two closed doors
			pick = random.choice([pick,remaining])
		elif strategy == 2:
			#strategy 2: Always switch
			pick = remaining 

		if winning == pick:
			wins_queue[strategy].put(1)
		plays_queue[strategy].put(1)
	

threads = []
for i in range(3):
	wins_queue.append(queue.Queue())
	plays_queue.append(queue.Queue())
	t = threading.Thread(target = create_sub_threads, args=[i])
	t.start()
	threads.append(t)	
for thread in threads:
	thread.join()