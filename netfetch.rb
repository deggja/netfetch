class Netfetch < Formula
  desc "CLI tool to scan for network policies in Kubernetes clusters/namespaces and provide a score based on the amount of untargeted workloads"
  homepage "https://github.com/deggja/netfetch"

  if OS.mac?
    url "https://github.com/deggja/netfetch/releases/download/5.2.4/netfetch_5.2.4_darwin_amd64.tar.gz"
    sha256 "364a20efe6533cfc06784f27e710885d3ccc0dd79abab463ed158ee6daf7066d"
  elsif OS.linux?
    url "https://github.com/deggja/netfetch/releases/download/5.2.4/netfetch_5.2.4_linux_amd64.tar.gz"
    sha256 "4982d0972fcb1b725b53c18df42e89eade7d558b2ea40e2da62407df42d537e1"
  end

  def install
    bin.install "netfetch"
  end

  test do
    system "#{bin}/netfetch", "version"
  end
end
