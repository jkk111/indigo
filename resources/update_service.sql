UPDATE services SET
  name = :name, desc = :desc, host = :host, path = :path, 
  repo = :repo, branch = :branch, hash = :hash, start = :start,
  args = :args, env = :env, install = :install, installArgs = :installArgs,
  installEnv = :installEnv, enabled = :enabled
WHERE id = :id