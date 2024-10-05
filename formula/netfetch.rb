class Netfetch < Formula
  desc "CLI tool to scan for network policies in Kubernetes clusters/namespaces and provide a score based on the amount of untargeted workloads"
  homepage "https://github.com/deggja/netfetch"

  if OS.mac?
    url "https://github.com/deggja/netfetch/releases/download/v0.5.3/netfetch_0.5.3_darwin_amd64.tar.gz"
    sha256 "96616cf0555a265289d4736fc09466c1e78fc9a13d6314a75c99eda81a9f7308"
  elsif OS.linux?
    url "https://github.com/deggja/netfetch/releases/download/v0.5.3/netfetch_0.5.3_linux_amd64.tar.gz"
    sha256 "96616cf0555a265289d4736fc09466c1e78fc9a13d6314a75c99eda81a9f7308"
  end

  def install
    bin.install "netfetch"
  end

  test do
    system "#{bin}/netfetch", "version"
  end
end
