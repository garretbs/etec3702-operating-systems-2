import random 
 
class MontyHall:
	def __init__(self):
		self.reset_doors()
		
	def reset_doors(self):
		self.doors = [False] * 3
		self.doors[self.choose_random_door()] = True
		
	def choose_random_door(self):
		return random.randint(0, 2)
 
	def select_door(self):
		self.selected_door = self.choose_random_door()
 
	def open_door(self):
		door_to_open = self.choose_random_door()
		while door_to_open == self.selected_door or self.doors[door_to_open]:
			door_to_open = self.choose_random_door()
		self.opened_door = door_to_open
 
	def switch_door(self):
		self.selected_door = 3 - self.selected_door - self.opened_door
		
	def stay_with_door(self):
		self.reset_doors()
		self.select_door()
		self.open_door()
		return self.doors[self.selected_door]
		
	def always_switch(self):
		self.reset_doors()
		self.select_door()
		self.open_door()
		self.switch_door()
		return self.doors[self.selected_door]
		
	def sometimes_switch(self):
		self.reset_doors()
		self.select_door()
		self.open_door()
		if random.randint(0, 1) == 1:
			self.switch_door()
		return self.doors[self.selected_door]
 
total_plays = 10000
stay_wins = 0
some_wins = 0
switch_wins = 0
monty = MontyHall()
for i in range(total_plays):
	stay_wins += 1 if monty.stay_with_door() else 0
	some_wins += 1 if monty.sometimes_switch() else 0
	switch_wins += 1 if monty.always_switch() else 0
print("Keep choice: \t\t", 100 * stay_wins / total_plays, "%")
print("Sometimes change: \t", 100 * some_wins / total_plays, "%")
print("Always switch: \t\t", 100 * switch_wins / total_plays, "%")