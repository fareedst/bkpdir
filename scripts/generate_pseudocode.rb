#!/usr/bin/env ruby
# frozen_string_literal: true

require 'yaml'
require 'fileutils'
require 'date'

TIED_DIR = File.join(__dir__, '..', 'tied')
IMPL_DIR = File.join(TIED_DIR, 'implementation-decisions')

def extract_refs(cross_refs, prefix)
  return [] unless cross_refs.is_a?(Array)
  cross_refs.select { |r| r.start_with?("#{prefix}-") }
end

def token_header(impl_token, arch_refs, req_refs, summary)
  tokens = [impl_token] + arch_refs + req_refs
  "# #{tokens.map { |t| "[#{t}]" }.join(' ')} — #{summary}"
end

def to_fn_name(text)
  text.to_s
    .gsub(/[`'"{}()\[\]]/, '')
    .gsub(/[^a-zA-Z0-9]/, '_')
    .gsub(/_+/, '_')
    .gsub(/^_|_$/, '')
    .downcase
    .slice(0, 50)
end

def step_to_pseudocode(detail_str, arch_refs, req_refs, impl_token)
  s = detail_str.to_s.strip
  return nil if s.empty?

  tokens_str = ([impl_token] + arch_refs + req_refs).map { |t| "[#{t}]" }.join(' ')

  if s.include?(':')
    label, rest = s.split(':', 2).map(&:strip)
    fn = to_fn_name(label)
    fn = 'step' if fn.empty?
    "  # #{tokens_str} #{s}\n  #{fn}(#{rest ? '"' + rest.slice(0, 60) + '"' : ''})"
  else
    fn = to_fn_name(s.split(/\s+/).first(4).join('_'))
    fn = 'step' if fn.empty?
    "  # #{tokens_str} #{s}\n  #{fn}()"
  end
end

def generate_pseudocode(token, data)
  record = data[token] || data
  name = record['name'] || token
  cross_refs = record['cross_references'] || []
  arch_refs = extract_refs(cross_refs, 'ARCH')
  req_refs = extract_refs(cross_refs, 'REQ')

  approach = record['implementation_approach'] || {}
  summary_text = (approach['summary'] || record['decision'] || name).to_s.split("\n").first.to_s.strip
  details = approach['details'] || []
  details = [details] unless details.is_a?(Array)
  details = details.uniq.reject { |d| d.to_s.strip.empty? }

  lines = []
  lines << token_header(token, arch_refs, req_refs, summary_text)
  lines << ""

  fn_name = to_fn_name(name)
  fn_name = 'implement' if fn_name.empty?

  if details.length > 0
    lines << "procedure #{fn_name}():"
    seen = {}
    details.each do |detail|
      d = detail.to_s.strip
      next if d.empty? || seen[d]
      seen[d] = true
      pseudo = step_to_pseudocode(d, arch_refs, req_refs, token)
      lines << pseudo if pseudo
    end
    lines << "  end"
  else
    lines << "procedure #{fn_name}():"
    rationale = record['rationale'] || {}
    why = rationale['why'].to_s.split("\n").first.to_s.strip
    lines << "  # #{([token] + arch_refs + req_refs).map { |t| "[#{t}]" }.join(' ')} #{why}" unless why.empty?
    lines << "  execute()"
    lines << "  end"
  end

  lines.join("\n")
end

def process_all_impl_files
  files = Dir.glob(File.join(IMPL_DIR, 'IMPL-*.yaml')).sort
  puts "Processing #{files.length} IMPL detail files..."

  success = 0
  errors = []

  files.each do |filepath|
    token = File.basename(filepath, '.yaml')
    begin
      data = YAML.safe_load(File.read(filepath), permitted_classes: [Date, Time])

      pseudocode = generate_pseudocode(token, data)

      tmpfile = "/tmp/pseudocode_#{token}.txt"
      File.write(tmpfile, pseudocode)

      yq_cmd = "yq -i '.\"#{token}\".essence_pseudocode = load_str(\"#{tmpfile}\")' '#{filepath}'"
      result = system(yq_cmd)

      if result
        system("yq -i -P '#{filepath}'")
        puts "  OK: #{token}"
        success += 1
      else
        errors << "#{token}: yq command failed"
        puts "  FAIL: #{token}"
      end

      File.delete(tmpfile) if File.exist?(tmpfile)
    rescue => e
      errors << "#{token}: #{e.message}"
      puts "  ERROR: #{token}: #{e.message}"
    end
  end

  puts ""
  puts "Results: #{success}/#{files.length} succeeded, #{errors.length} errors"
  errors.each { |e| puts "  - #{e}" } unless errors.empty?
end

process_all_impl_files
