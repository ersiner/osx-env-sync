#!/usr/bin/env ruby

## Pick Bash or Zsh:
DEFAULT_SHELL = 'zsh'

class Command
    CopiedEnvs = %w[HOME LOGNAME USER LANG]

    attr_reader :default_shell
    def initialize(shell: nil, prefix: nil, suffix: [])
        @shell = shell
        @prefix = prefix || clean_shell_env
        @suffix = suffix
        @default_shell = DEFAULT_SHELL
    end

    def shell
        @shell || env_shell || default_shell
    end

    def env_shell
        return nil unless ENV['SHELL'] && !ENV['SHELL'].empty?
        ENV['SHELL']
    end

    def mapped_envs
        CopiedEnvs.map {|e| "#{e}='#{ENV[e]}'" }
    end

    def clean_shell_env
        ["env -i"] + mapped_envs + ["TERM=xterm"]
    end

    def shell_command
        case shell
        when /bash/
            return shell + " --login -c env"

        when /zsh/
            return shell + " -i --login -c env"

        else
            $stderr.puts "I don't know how to work with shell #{shell.inspect}."
            exit 1
        end
    end

    def command
        return (@prefix + [shell_command] + @suffix).join ' '
    end

    def run
        cmd = command
        $stderr.puts "# Running: #{cmd}"
        output = `#{cmd}`.split
        $stderr.puts "Output = ", output.inspect
        return output
    end
end

def run_setenv(name, val)
    cmd = ['launchctl', 'setenv', name, val]
    puts "# Running: #{cmd.inspect}"
    system *cmd
end

Command.new.run.collect {|i| /^(.+)=(.+)$/.match i}.select {|i| i}.collect {|i| i[1..2]}.each {|i| run_setenv(*i)}
