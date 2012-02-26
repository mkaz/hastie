#!/usr/bin/env ruby

# determine arch
if ENV['GOARCH'] == "386" then
  compile_cmd = "8g"
  link_cmd = "8l"
else
  compile_cmd = "6g"
  link_cmd = "6l"
end


task :build do

  filename = "hastie"

  # collect stats on build
  # TODO

  # compile
  system "#{compile_cmd} -o #{filename}.link #{filename}.go"

  ## verify no error
  system "#{link_cmd} -o #{filename} #{filename}.link"

  ## verify link exists
  system "rm #{filename}.link"
end


task :test do
  system "cd test; ../hastie"
end

task :verify do
  system "cd test; cat public/index.html"
end

