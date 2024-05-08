class Netfetch < Formula
  desc "CLI tool to scan for network policies in Kubernetes clusters/namespaces and provide a score based on the amount of untargeted workloads"
  homepage "https://github.com/deggja/netfetch"

  if OS.mac?
    url "https://github.com/deggja/netfetch/releases/download/5.0.1/netfetch_5.0.1_darwin_amd64.tar.gz"
    sha256 "35674d84e214bd096339fec2a2c27499d073531a219ab7d3f56a6ef5b874ba8d"
  elsif OS.linux?
    url "https://github.com/deggja/netfetch/releases/download/5.0.1/netfetch_5.0.1_linux_amd64.tar.gz"
    sha256 "442f571b92376429b56de827a6e00628afbd59a20f932fc6fe566a2440ed2877"
  end

  def install
    bin.install "netfetch"
  end

  test do
    system "#{bin}/netfetch", "version"
  end
end
