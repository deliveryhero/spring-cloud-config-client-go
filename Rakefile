task :command_exists, [:command] do |_, args|
    abort "#{args.command} doesn't exists" if `command -v #{args.command} > /dev/null 2>&1 && echo $?`.chomp.empty?
end

task :is_repo_clean do
  abort 'please commit your changes first!' unless `git status -s | wc -l`.strip.to_i.zero?
end

task :current_branch do
  `git rev-parse --abbrev-ref HEAD`.strip
end

task :has_bumpversion do
  Rake::Task['command_exists'].invoke('bumpversion')
end

task :has_hadolint do
    Rake::Task['command_exists'].invoke('hadolint')
end

task :has_golangcilint do
    Rake::Task['command_exists'].invoke('golangci-lint')
end
  
task :has_pre_commit do
    Rake::Task['command_exists'].invoke('pre-commit')
end

task :has_gsed do
    Rake::Task['command_exists'].invoke('gsed')
end

desc "run golangci-lint"
task :lint => [:has_golangcilint] do
  system "LOG_LEVEL=error golangci-lint run"
end

namespace :pre_commit do
    desc "update hooks"
    task :update => [:has_pre_commit] do
      system "pre-commit autoupdate"
    end
end

desc "publish new version of the library, default is: patch"
task :publish, [:revision] => [:is_repo_clean] do |t, args|
  current_branch = Rake::Task["current_branch"].invoke.first.call

  abort "this command works for [main] or [master] branches only, not for [#{current_branch}] branch" unless ['master', 'main'].include?(current_branch)

  args.with_defaults(revision: "patch")

  Rake::Task["bump"].invoke(args.revision)
  current_git_tag = "v#{current_version}"

  puts "[->] new version is \e[33m#{current_git_tag}\e[0m"
  puts "[->] pushing \e[33m#{current_git_tag}\e[0m to remote"
  system %{
    git push origin #{current_git_tag} &&
    go list -m github.com:deliveryhero/sc-honeylogger@#{current_git_tag} &&
    echo "[->] [#{current_git_tag}] has been published" &&
    git push origin #{current_branch} &&
    echo "[->] code pushed to: [#{current_branch}] branch (updated)"
  }
end

AVAILABLE_REVISIONS = %w[major minor patch].freeze
desc "bump version, default is: patch"
task :bump, [:revision] => [:has_bumpversion] do |_, args|
  args.with_defaults(revision: 'patch')
  unless AVAILABLE_REVISIONS.include?(args.revision)
    abort "Please provide valid revision: #{AVAILABLE_REVISIONS.join(',')}"
  end

  system "bumpversion #{args.revision}"
end

desc "run tests"
task :test => [:has_gsed] do
  system %{
    color_red=$'\e[0;31m'
    color_yellow=$'\e[0;33m'
    color_white=$'\e[0;37m'
    color_off=$'\e[0m'
    
    any_errors="0"
    
    for s in $(go list ./...); do 
      if ! go test -failfast -p 1 -v -race "${s}"; then
        echo "\t\n${color_red}${s}${color_off} ${color_yellow}fails${color_off}...\n"
        any_errors="1"
        break
      fi
    done
    
    if [[ "${any_errors}" == "0" ]]; then
      echo "\n\n${color_white}Tests are passing...${color_off}"
      echo "${color_white}Calculating code coverage${color_off}"
      go test ./... -coverpkg=./src/... -coverprofile ./coverage.out > /dev/null 2>&1
      code_coverage_ratio=$(go tool cover -func ./coverage.out | grep "total:" | awk '{print $3}')
      echo "${color_white}Total test coverage: ${color_yellow}${code_coverage_ratio}${color_off}"
      code_coverage_ratio_md=${code_coverage_ratio/%/25}
      gsed -i -r "s/coverage-[0-9\.\%]+/coverage-${code_coverage_ratio_md}/" README.md &&
      echo "README updated...\n"
    fi
  }
end

task :update_test_coverage => [:has_gsed] do
  system %{
    color_red=$'\e[0;31m'
    color_yellow=$'\e[0;33m'
    color_white=$'\e[0;37m'
    color_off=$'\e[0m'

    echo "${color_white}Calculating code coverage${color_off}"
    go test ./... -coverpkg=./src/... -coverprofile ./coverage.out > /dev/null 2>&1
    code_coverage_ratio=$(go tool cover -func ./coverage.out | grep "total:" | awk '{print $3}')
    echo "${color_white}Total test coverage: ${color_yellow}${code_coverage_ratio}${color_off}"
    code_coverage_ratio_md=${code_coverage_ratio/%/25}
    gsed -i -r "s/coverage-[0-9\.\%]+/coverage-${code_coverage_ratio_md}/" README.md &&
    echo "README updated...\n"
  }
end