#!/usr/bin/env ruby
# frozen_string_literal: true

# Ruby 3.0.3 script to convert semantic token format from colon to hyphen
# in non-Go files: shell scripts, Makefile, markdown, YAML, etc.
#
# Usage: ruby scripts/convert_tokens_nongo.rb [--dry-run]

require 'find'

DRY_RUN = ARGV.include?('--dry-run')
ROOT = File.expand_path('..', __dir__)

TOKEN_PREFIXES = %w[REQ ARCH IMPL TEST PROC ACTION].freeze
BRACKET_RE = /\[(#{TOKEN_PREFIXES.join('|')}):([A-Za-z0-9_]+)\]/

EXTENSIONS = %w[.sh .md .py .yaml .yml .txt].freeze
BASENAMES = %w[Makefile].freeze

SKIP_DIRS = %w[vendor .git node_modules tied].freeze
SKIP_PATTERNS = [/\.plan\.md$/].freeze

def target_files
  files = []
  Find.find(ROOT) do |path|
    Find.prune if SKIP_DIRS.any? { |d| path.include?("/#{d}/") || File.basename(path) == d }
    next unless File.file?(path)
    next if SKIP_PATTERNS.any? { |p| path.match?(p) }

    ext = File.extname(path)
    base = File.basename(path)
    files << path if EXTENSIONS.include?(ext) || BASENAMES.include?(base)
  end
  files.sort
end

def convert_line(line)
  line.gsub(BRACKET_RE) { "[#{$1}-#{$2}]" }
end

def process_file(filepath)
  original = File.read(filepath, encoding: 'UTF-8')
  unless original.valid_encoding?
    original = original.encode('UTF-8', 'binary', invalid: :replace, undef: :replace, replace: '')
  end
  lines = original.lines
  changed_lines = []

  new_lines = lines.map.with_index do |line, idx|
    converted = convert_line(line)
    if converted != line
      changed_lines << { line_num: idx + 1, before: line.rstrip, after: converted.rstrip }
    end
    converted
  end

  result = new_lines.join
  return nil if result == original

  unless DRY_RUN
    File.write(filepath, result)
  end

  { file: filepath.sub("#{ROOT}/", ''), changes: changed_lines }
end

puts DRY_RUN ? "DRY RUN - no files will be modified" : "LIVE RUN - files will be modified"
puts "(Skipping tied/ directory - those files already use hyphen format)"
puts ""

files = target_files
puts "Scanning #{files.length} non-Go files..."
puts ""

total_files_changed = 0
total_lines_changed = 0

files.each do |filepath|
  result = process_file(filepath)
  next unless result

  total_files_changed += 1
  total_lines_changed += result[:changes].length

  puts "#{result[:file]} (#{result[:changes].length} changes)"
  result[:changes].first(3).each do |c|
    puts "  L#{c[:line_num]}: #{c[:before][0..100]}"
    puts "     -> #{c[:after][0..100]}"
  end
  if result[:changes].length > 3
    puts "  ... and #{result[:changes].length - 3} more changes"
  end
  puts ""
end

puts "=" * 60
puts "Summary: #{total_files_changed} files, #{total_lines_changed} lines changed"
puts DRY_RUN ? "(dry run - no changes written)" : "(changes applied)"
