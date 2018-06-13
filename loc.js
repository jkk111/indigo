let fs = require('fs')
let path = require('path')
console.log(__filename)

let get_files = (dir) => {
  let _files = fs.readdirSync(dir);
  let files = [];
  for(var file of _files) {
    let filepath = path.join(dir, file)
    let stat = fs.statSync(filepath);

    if(stat.isDirectory()) {
      files.push(...get_files(filepath))
    } else {
      files.push(filepath)
    }
  }

  return files;
}

let loc = 0;

let files = get_files(__dirname)
files.splice(files.indexOf(path.join(__dirname, 'assets', 'bindata.go')), 1)
files = files.filter(f => f.indexOf(path.join(__dirname, '.git')) != 0)
files.splice(files.indexOf(path.join(__dirname, 'indigo.exe')), 1)

for(var file of files) {
  let data = fs.readFileSync(file, 'utf8')
  loc += data.split('\n').length
}

console.log(loc)