class Netfetch < Formula
    desc "CLI tool to scan for network policies in Kubernetes clusters/namespaces and provide a score based on the amount of untargeted workloads"
    homepage "https://github.com/deggja/netfetch"
    url "https://github.com/deggja/netfetch/releases/download/1.0.0/netfetch_1.0.0_darwin_amd64.tar.gz"
    sha256 "9e9937765e2c0b11b94a1d82f65b623fd7fc3504b99916fdb6d47449ae020b83"
  
    def install
      bin.install "netfetch"
    end
  
    test do
      system "#{bin}/netfetch", "--version"
    end
  end