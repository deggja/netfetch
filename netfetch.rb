class Netfetch < Formula
    desc "CLI tool to scan for network policies in Kubernetes clusters/namespaces and provide a score based on the amount of untargeted workloads"
    homepage "https://github.com/deggja/netfetch"
    url "https://github.com/deggja/netfetch/releases/download/1.0.4/netfetch_1.0.4_darwin_amd64.tar.gz"
    sha256 "c74157ff160837d0878db010f23aa5e506e4463308448e439023fad9991d2489"
  
    def install
      bin.install "netfetch"
    end
  
    test do
      system "#{bin}/netfetch", "--version"
    end
  end