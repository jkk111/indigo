INSERT INTO services (
  name, desc, host, path, 
  repo, branch, hash, start,
  args, env, install, installArgs,
  installEnv, enabled
) VALUES (
  :name, :desc, :host, :path, 
  :repo, :branch, :hash, :start,
  :args, :env, :install, :installArgs,
  :installEnv, :enabled
)