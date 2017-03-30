import urllib.request
import re
import concurrent.futures
import threading
from html.parser import HTMLParser

class HP(HTMLParser):
	def __init__(self, site):
		super().__init__()
		self.current_site = site
		self.result = []
		
	def handle_starttag(self, tag, attributes): #<a> <b> <i> etc.
		global sites_to_crawl
		if tag == 'a':
			for attr in attributes:
				if attr[0] == 'href':
					if attr[1][0] != '#': #same page 
						new_site = urllib.parse.urljoin(self.current_site, attr[1])
						domain = urllib.parse.urlparse(new_site).netloc
						if domain == main_domain: #same hostname as input
							with lock:
								if new_site not in visited:
									self.result.append(new_site)
								
		elif tag == 'meta':
			yoyo = is_keyword.match(str(attributes[0][1]))
			if yoyo:
				with lock:
					words_list = is_a_word.findall(attributes[1][1])
					for word in words_list:
						if word not in words:
							words[word] = []
						if self.current_site not in words[word]:
							words[word].append(self.current_site)
				
					
	def handle_endtag(self, tag): #</a> </b> </i> etc.
		pass
	
	def handle_data(self, text): #things between tags
		with lock:
			words_list = is_a_word.findall(text)
			for word in words_list:
				if word not in words:
					words[word] = []
				if self.current_site not in words[word]:
					words[word].append(self.current_site)

def crawl_site(site):
	parser = HP(site)
	try:
		url = urllib.request.urlopen(site)
	except:
		parser.close()
		return False
	d = url.read()
	try:
		raw_html = d.decode()
	except:
		raw_html = d.decode("latin-1")
	url.close()
	parser.feed(raw_html)
	for res in parser.result:
		with lock:
			if res not in visited:
				visited.append(res)
				futures.append(pool.submit(crawl_site, res))
	parser.close()
	return True

print("Enter a site to spider: ")
main_site = input()
main_domain = urllib.parse.urlparse(main_site).netloc
print()

#is_a_word = re.compile('([a-zA-Z_]+)')
is_a_word = re.compile('(\w+)')
is_keyword = re.compile('keywords', re.IGNORECASE)

num_threads = 4
pool = concurrent.futures.ThreadPoolExecutor(num_threads)
lock = threading.RLock() #globals lock

words = {}
visited = [main_site]
futures = [pool.submit(crawl_site, main_site)]

for future in futures:
	future.result()


for word in words:
	print(word, "=>")
	for site in words[word]:
		print(site)
	print()

