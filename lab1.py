import random
import threading



def monty_hall(strategy):
	win=0
	total_plays = 10000
	for trials in range(total_plays):
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
			win+=1
		
	print("Strategy",strategy,win/total_plays*100,"%")

threads = []
for i in range(3):
	t = threading.Thread(target = monty_hall, args=[i])
	t.start()
	threads.append(t)
for thread in threads:
	thread.join()