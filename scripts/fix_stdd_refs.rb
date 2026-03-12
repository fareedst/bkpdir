#!/usr/bin/env ruby
# frozen_string_literal: true

# Ruby 3.0.3 script to update stdd/ path references to tied/ in tied/*.md files.
# Also updates stdd/ references in root-level files.

require 'find'

DRY_RUN = ARGV.include?('--dry-run')
ROOT = File.expand_path('..', __dir__)

REPLACEMENTS = {
  'stdd/requirements.md' => 'tied/requirements.md',
  'stdd/architecture-decisions.md' => 'tied/architecture-decisions.md',
  'stdd/implementation-decisions.md' => 'tied/implementation-decisions.md',
  'stdd/semantic-tokens.md' => 'tied/semantic-tokens.md',
  'stdd/ai-principles.md' => 'tied/ai-principles.md',
  'stdd/ai-assistant-compliance.md' => 'tied/ai-assistant-compliance.md',
  'stdd/merge-strategy-conflict-resolution.md' => 'tied/merge-strategy-conflict-resolution.md',
  'stdd/tasks.md' => 'tied/tasks.md',
  '`stdd/' => '`tied/',
  'stdd/*.md' => 'tied/*.md',
}

STDD_PATH_RE = /\bstdd\//

def target_files
  files = []
  # tied/ markdown files
  Dir.glob(File.join(ROOT, 'tied', '*.md')).each { |f| files << f }
  # Root-level files
  %w[README.md CHANGELOG.md .ai-agent-instructions ai-principles.md].each do |name|
    path = File.join(ROOT, name)
    files << path if File.exist?(path)
  end
  files.sort.uniq
end

def process_file(filepath)
  original = File.read(filepath, encoding: 'UTF-8')
  modified = original.dup

  REPLACEMENTS.each do |old, new_val|
    modified.gsub!(old, new_val)
  end

  return nil if modified == original

  changes = 0
  original.lines.zip(modified.lines).each_with_index do |(o, m), idx|
    changes += 1 if o != m
  end

  unless DRY_RUN
    File.write(filepath, modified)
  end

  { file: filepath.sub("#{ROOT}/", ''), changes: changes }
end

puts DRY_RUN ? "DRY RUN" : "LIVE RUN"
puts ""

files = target_files
total_files = 0
total_changes = 0

files.each do |f|
  result = process_file(f)
  next unless result

  total_files += 1
  total_changes += result[:changes]
  puts "#{result[:file]} (#{result[:changes]} lines)"
end

puts ""
puts "Summary: #{total_files} files, #{total_changes} lines changed"
