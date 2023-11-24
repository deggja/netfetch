class Netfetch < Formula
    desc "CLI tool to scan for network policies in Kubernetes clusters"
    homepage "https://github.com/deggja/netfetch"
<<<<<<< HEAD
    url "https://github.com/deggja/netfetch/releases/download/0.0.15/netfetch_0.0.15_darwin_amd64.tar.gz"
    sha256 "2bbf7cb3d50e477d45d6557417048a07b1d1f3884543cad3583049ef04072857"
=======
    url "https://github.com/deggja/netfetch/releases/download/0.0.9/netfetch_0.0.9_darwin_amd64.tar.gz"
    sha256 "53bfbc3240986e4774ec2e7652be4bcac58c05e118d5443cbf0dd37ff8572722"
>>>>>>> be953cb5bea308074e3d9979057f2d9e94e5214e
  
    def install
      bin.install "netfetch"
    end
  
    test do
      system "#{bin}/netfetch", "--version"
    end
  end
<<<<<<< HEAD
=======
  
>>>>>>> be953cb5bea308074e3d9979057f2d9e94e5214e
