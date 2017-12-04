const formFieldName = 'configs';
const validFilenames = ['server_cfg.ini', 'entry_list.ini'];

// Code partially used:
// Copyright (c) silverwind https://github.com/silverwind/uppie

export default function handleDataTransfer(e, cb) {
  const dt = e.dataTransfer;
  if (dt.items && dt.items.length && 'webkitGetAsEntry' in dt.items[0]) {
    entriesApi(dt.items, cb.bind(null, e));
  } else if ('getFilesAndDirectories' in dt) {
    newDirectoryApi(dt, cb.bind(null, e));
  } else if (dt.files) {
    arrayApi(dt, cb.bind(null, event));
  } else {
    cb();
  }
}

// API implemented in Firefox 42+ and Edge
function newDirectoryApi(input, cb) {
  const fd = new FormData(), files = [];
  const iterate = function(entries, path, resolve) {
    const promises = [];
    entries.forEach(function(entry) {
      promises.push(new Promise(function(resolve) {
        if ('getFilesAndDirectories' in entry) {
          entry.getFilesAndDirectories().then(function(entries) {
            iterate(entries, entry.path + '/', resolve);
          });
        } else {
          if (entry.name && validFilenames.indexOf(entry.name) !== -1) {
            const p = (path + entry.name).replace(/^[/\\]/, '');
            fd.append(formFieldName, entry, p);
            files.push(p);
          }
          resolve();
        }
      }));
    });
    Promise.all(promises).then(resolve);
  };
  input.getFilesAndDirectories().then(function(entries) {
    new Promise(function(resolve) {
      iterate(entries, '/', resolve);
    }).then(cb.bind(null, fd, files));
  });
}

// old prefixed API implemented in Chrome 11+ as well as array fallback
function arrayApi(input, cb) {
  const fd = new FormData(), files = [];
  [].slice.call(input.files).forEach(function(file) {
    if (validFilenames.indexOf(file.name) !== -1) {
      fd.append(formFieldName, file, file.webkitRelativePath || file.name);
      files.push(file.webkitRelativePath || file.name);
    }
  });
  cb(fd, files);
}

// old drag and drop API implemented in Chrome 11+
function entriesApi(items, cb) {
  const fd = new FormData(), files = [], rootPromises = [];

  function readEntries(entry, reader, oldEntries, cb) {
    const dirReader = reader || entry.createReader();
    dirReader.readEntries(function(entries) {
      const newEntries = oldEntries ? oldEntries.concat(entries) : entries;
      if (entries.length) {
        setTimeout(readEntries.bind(null, entry, dirReader, newEntries, cb), 0);
      } else {
        cb(newEntries);
      }
    });
  }

  function readDirectory(entry, path, resolve) {
    if (!path) path = entry.name;
    readEntries(entry, 0, 0, function(entries) {
      const promises = [];
      entries.forEach(function(entry) {
        promises.push(new Promise(function(resolve) {
          if (entry.isFile) {
            entry.file(function(file) {
              if (validFilenames.indexOf(file.name) !== -1) {
                const p = path + '/' + file.name;
                fd.append(formFieldName, file, p);
                files.push(p);
              }
              resolve();
            }, resolve.bind());
          } else readDirectory(entry, path + '/' + entry.name, resolve);
        }));
      });
      Promise.all(promises).then(resolve.bind());
    });
  }

  [].slice.call(items).forEach(function(entry) {
    entry = entry.webkitGetAsEntry();
    if (entry) {
      rootPromises.push(new Promise(function(resolve) {
        if (entry.isFile) {
          entry.file(function(file) {
            if (validFilenames.indexOf(file.name) !== -1) {
              fd.append(formFieldName, file, file.name);
              files.push(file.name);
            }
            resolve();
          }, resolve.bind());
        } else if (entry.isDirectory) {
          readDirectory(entry, null, resolve);
        }
      }));
    }
  });
  Promise.all(rootPromises).then(cb.bind(null, fd, files));
}
