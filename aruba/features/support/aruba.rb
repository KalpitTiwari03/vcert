require 'aruba/cucumber'
require "json_spec/cucumber"

Aruba.configure do |config|
  config.allow_absolute_paths = true
end

$path_separator = "/"

$temp_path = "tmp/aruba"

def last_json
  last_command_started.stdout.to_s
end

def random_cn
  Time.now.to_i.to_s + "-" + (0..4).to_a.map{|a| rand(36).to_s(36)}.join + ".venafi.example.com"
end

def random_string
  Time.now.to_i.to_s + "-" + (0..4).to_a.map{|a| rand(36).to_s(36)}.join
end

def random_filename
  Time.now.to_i.to_s + "-" + (0..6).to_a.map{|a| rand(36).to_s(36)}.join + ".txt"
end
