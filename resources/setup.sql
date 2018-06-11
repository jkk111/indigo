CREATE TABLE IF NOT EXISTS services (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  desc TEXT NOT NULL DEFAULT "echo 'No Description Specified'",
  host TEXT NOT NULL, 
  path TEXT NOT NULL, 
  repo TEXT NOT NULL DEFAULT "echo 'No Repo Specified'",
  start TEXT NOT NULL DEFAULT "echo 'No Start Command Specified'",
  args TEXT NOT NULL DEFAULT "[]",
  env TEXT NOT NULL DEFAULT "[]",
  install TEXT NOT NULL DEFAULT "echo 'No Install Command Specified'",
  installArgs TEXT NOT NULL DEFAULT "[]",
  installEnv TEXT NOT NULL DEFAULT "[]",
  enabled BOOL NOT NULL DEFAULT 1,
  UNIQUE(host, path)
  UNIQUE(name)
);

CREATE TABLE IF NOT EXISTS admins (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user TEXT NOT NULL,
  serial TEXT NOT NULL,
  UNIQUE(serial)
);