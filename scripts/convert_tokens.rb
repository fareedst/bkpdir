#!/usr/bin/env ruby
# frozen_string_literal: true

# Ruby 3.0.3 script to convert semantic token format from colon to hyphen
# in Go source and test files.
#
# Converts: [REQ:X] -> [REQ-X], [ARCH:X] -> [ARCH-X], [IMPL:X] -> [IMPL-X],
#           [TEST:X] -> [TEST-X], [PROC:X] -> [PROC-X]
#
# Also handles bare references in comments (REQ:X -> REQ-X) but avoids
# changing Go map syntax, YAML tags, or struct field assignments.
#
# Usage: ruby scripts/convert_tokens.rb [--dry-run]

require 'find'

DRY_RUN = ARGV.include?('--dry-run')
ROOT = File.expand_path('..', __dir__)

TOKEN_PREFIXES = %w[REQ ARCH IMPL TEST PROC ACTION].freeze
BRACKET_RE = /\[(#{TOKEN_PREFIXES.join('|')}):([A-Za-z0-9_]+)\]/
BARE_COMMENT_RE = /(?<=\s|^)(#{TOKEN_PREFIXES.join('|')}):([A-Za-z0-9_]+)(?=[\s,\]\)]|$)/

SKIP_DIRS = %w[vendor .git node_modules].freeze

def go_files
  files = []
  Find.find(ROOT) do |path|
    Find.prune if SKIP_DIRS.any? { |d| path.include?("/#{d}/") }
    files << path if path.end_with?('.go')
  end
  files.sort
end

def convert_line(line)
  # Skip lines that are Go code (not comments) containing map/struct syntax
  # Only convert tokens in comment lines (starting with //) or in string literals
  # that reference tokens

  modified = line.dup

  # Convert bracketed tokens: [REQ:X] -> [REQ-X]
  modified.gsub!(BRACKET_RE) { "[#{$1}-#{$2}]" }

  # Convert bare tokens in comments only (lines starting with //)
  if line.strip.start_with?('//')
    modified.gsub!(BARE_COMMENT_RE) { "#{$1}-#{$2}" }
  end

  modified
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
puts ""

files = go_files
puts "Scanning #{files.length} Go files..."
puts ""

total_files_changed = 0
total_lines_changed = 0

files.each do |filepath|
  result = process_file(filepath)
  next unless result

  total_files_changed += 1
  total_lines_changed += result[:changes].length

  puts "#{result[:file]} (#{result[:changes].length} changes)"
  result[:changes].first(5).each do |c|
    puts "  L#{c[:line_num]}: #{c[:before]}"
    puts "     -> #{c[:after]}"
  end
  if result[:changes].length > 5
    puts "  ... and #{result[:changes].length - 5} more changes"
  end
  puts ""
end

puts "=" * 60
puts "Summary: #{total_files_changed} files, #{total_lines_changed} lines changed"
puts DRY_RUN ? "(dry run - no changes written)" : "(changes applied)"
