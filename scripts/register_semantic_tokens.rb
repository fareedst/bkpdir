#!/usr/bin/env ruby
# frozen_string_literal: true

# Ruby 3.0.3 script to register all project-specific REQ/ARCH/IMPL tokens
# from the TIED indexes into semantic-tokens.yaml via yq.

require 'yaml'
require 'date'

TIED_DIR = File.join(__dir__, '..', 'tied')
ST_FILE = File.join(TIED_DIR, 'semantic-tokens.yaml')

def load_yaml(path)
  YAML.safe_load(File.read(path), permitted_classes: [Date, Time, Symbol]) || {}
rescue Psych::DisallowedClass => e
  # Fall back to YAML.load for complex types
  YAML.load(File.read(path), permitted_classes: :all) || {}
rescue => e
  warn "Warning: Failed to load #{path}: #{e.message}"
  {}
end

def existing_tokens
  data = load_yaml(ST_FILE)
  data.keys.reject { |k| k.start_with?('#') }
end

def source_index_for(type)
  case type
  when 'REQ' then 'requirements.yaml'
  when 'ARCH' then 'architecture-decisions.yaml'
  when 'IMPL' then 'implementation-decisions.yaml'
  end
end

def detail_dir_for(type)
  case type
  when 'REQ' then 'requirements'
  when 'ARCH' then 'architecture-decisions'
  when 'IMPL' then 'implementation-decisions'
  end
end

def type_name_for(prefix)
  prefix
end

existing = existing_tokens
puts "Existing tokens in semantic-tokens.yaml: #{existing.length}"

new_tokens = []

%w[requirements architecture-decisions implementation-decisions].each do |index_name|
  index_path = File.join(TIED_DIR, "#{index_name}.yaml")
  data = load_yaml(index_path)
  next unless data.is_a?(Hash)

  data.each do |token, record|
    next if existing.include?(token)
    next unless record.is_a?(Hash)

    prefix = token.split('-').first
    type = type_name_for(prefix)
    name = record['name'] || token
    status = record['status'] || 'Active'
    description = record['description'] || name
    cross_refs = record['cross_references'] || []
    detail_file = record['detail_file'] || "#{detail_dir_for(prefix)}/#{token}.yaml"
    category = record['category']

    new_tokens << {
      token: token,
      type: type,
      name: name,
      category: category,
      status: status,
      description: description,
      cross_references: cross_refs,
      source_index: source_index_for(prefix),
      detail_file: detail_file
    }
  end
end

puts "New tokens to register: #{new_tokens.length}"

if new_tokens.empty?
  puts "Nothing to do."
  exit 0
end

today = Date.today.strftime('%Y-%m-%d')

new_tokens.each do |t|
  tmpfile = "/tmp/st_#{t[:token]}.yaml"

  yaml_content = {
    t[:token] => {
      'type' => t[:type],
      'name' => t[:name],
      'status' => t[:status] == 'Implemented' ? 'Active' : t[:status],
      'description' => t[:description],
      'cross_references' => t[:cross_references],
      'source_index' => t[:source_index],
      'detail_file' => t[:detail_file],
      'metadata' => {
        'registered' => today,
        'last_updated' => today
      }
    }
  }
  yaml_content[t[:token]]['category'] = t[:category] if t[:category]

  File.write(tmpfile, YAML.dump(yaml_content).sub(/^---\n/, ''))

  cmd = "yq -i '. *= load(\"#{tmpfile}\")' '#{ST_FILE}'"
  result = system(cmd)

  if result
    puts "  OK: #{t[:token]}"
  else
    puts "  FAIL: #{t[:token]}"
  end

  File.delete(tmpfile) if File.exist?(tmpfile)
end

system("yq -i -P '#{ST_FILE}'")
puts ""
puts "Done. Validating..."

final_data = load_yaml(ST_FILE)
final_count = final_data.keys.reject { |k| k.start_with?('#') }.length
puts "Total tokens in semantic-tokens.yaml: #{final_count}"
