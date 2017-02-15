import os
import shutil

if os.path.exists("release") == False:
	os.mkdir("release")

print("[OK] bulid release tool start...")

print("[OK] compiling go ...")

code = os.system("GOOS=linux GOARCH=amd64 go build -o ./release/collector_linux64 main.go")
if code > 0:
	raise ValueError("build error")

dirSrc = "./"
dirDest = "./release"

def getDestPath(path, src, dest):
	if src == "./":
		if path.find("./") == 0:
			path = path[2:]
		return os.path.join(dest, path)

	return path.replace(src, dest)

def hasHidePath(path):
	paths = path.split(os.sep)
	if paths[0] == ".":
		paths = paths[1:]

	if len(paths[0]) <= 0:
		return False

	for p in paths:
		if p[0] == ".":
			return True

	return False



# set files and paths that should be exclued from release files
ignoredDirs = ["parser", "data", "release"]
ignoredFiles = ["build_release.py", "main.go", "README.md", "conf/config.json"]


ignoredPaths = []

for d in ignoredDirs:
	ignoredPaths.append(os.path.join(dirSrc, d))

for f in ignoredFiles:
	ignoredPaths.append(os.path.join(dirSrc, f))

def sha1OfFile(filepath):
    import hashlib
    sha = hashlib.sha1()
    with open(filepath, 'rb') as f:
        while True:
            block = f.read(2**10) # Magic number: one-megabyte blocks.
            if not block: break
            sha.update(block)
        return sha.hexdigest()

def isPathIgnored(filepath):
	for p in ignoredPaths:
		index = filepath.find(p) 
		if index == 0:
			return True

	return False


def buildRelease():

	for root, dirs, files in os.walk(dirSrc):
		if hasHidePath(root):
			continue

		for x in dirs:
			srcPath = os.path.join(root,x)
			destRelativePath = getDestPath(srcPath, dirSrc, dirDest)

			if hasHidePath(x):
				continue

			if isPathIgnored(srcPath) == True:
				continue

			if os.path.exists(destRelativePath) == False:
				print(destRelativePath, " dir does not exists")
				try:
					os.mkdir(destRelativePath)
				except OSError as e:
					print("[ERR] mkdir error: ", e)
					raise e
		

		for x in files:
			filePath = os.path.join(root,x)
			destRelativePath = getDestPath(filePath, dirSrc, dirDest)

			if hasHidePath(x):
				continue

			if isPathIgnored(filePath) == True:
				continue

			if os.path.exists(destRelativePath) == False:
				try:
					shutil.copy(filePath, destRelativePath)
				except Exception as e:
					print("[ERR] copy file error: ", filePath, " ", e)
					raise e
			else:
				destSha1 = sha1OfFile(destRelativePath)
				srcSha1 = sha1OfFile(filePath)
				if destSha1 != srcSha1:
					try:
						os.remove(destRelativePath)
						shutil.copy(filePath, destRelativePath)
					except Exception as e:
						raise e


	print("--------------------------------------------")
	print("build OK")

print("[OK] copying files ...")
buildRelease()
print("[OK] push binary file to git ...")


code = os.system("cd release")
if code > 0:
	raise ValueError("cd error")

code = os.system("git add --all")
if code > 0:
	raise ValueError("git add error")

code = os.system('git commit -m "new build" ')
if code > 0:
	raise ValueError("commit error")

code = os.system('git push ')
if code > 0:
	raise ValueError("push error")
	
print("[OK] all done")
