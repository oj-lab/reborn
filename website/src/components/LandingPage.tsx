import { useEffect } from 'react'

export default function LandingPage() {
  useEffect(() => {
    // Add custom CSS animations
    const style = document.createElement('style')
    style.textContent = `
      @keyframes fadeInSlideDown {
        from {
          opacity: 0;
          transform: translateY(-16px);
        }
        to {
          opacity: 1;
          transform: translateY(0);
        }
      }
      
      @keyframes fadeInSlideUp {
        from {
          opacity: 0;
          transform: translateY(32px);
        }
        to {
          opacity: 1;
          transform: translateY(0);
        }
      }
      
      @keyframes scaleX {
        from {
          transform: scaleX(0);
        }
        to {
          transform: scaleX(1);
        }
      }
      
      @keyframes floatLine1 {
        0% {
          transform: translateX(-20px) translateY(-15px);
          opacity: 0.2;
        }
        25% {
          transform: translateX(0px) translateY(-5px);
          opacity: 0.5;
        }
        50% {
          transform: translateX(20px) translateY(10px);
          opacity: 0.8;
        }
        75% {
          transform: translateX(0px) translateY(-5px);
          opacity: 0.5;
        }
        100% {
          transform: translateX(-20px) translateY(-15px);
          opacity: 0.2;
        }
      }
      
      @keyframes floatLine2 {
        0% {
          transform: translateX(15px) translateY(-10px) scale(0.9);
          opacity: 0.3;
        }
        25% {
          transform: translateX(0px) translateY(2px) scale(1);
          opacity: 0.6;
        }
        50% {
          transform: translateX(-15px) translateY(15px) scale(1.1);
          opacity: 0.9;
        }
        75% {
          transform: translateX(0px) translateY(2px) scale(1);
          opacity: 0.6;
        }
        100% {
          transform: translateX(15px) translateY(-10px) scale(0.9);
          opacity: 0.3;
        }
      }
      
      @keyframes floatLine3 {
        0% {
          transform: translateX(-10px) translateY(20px) scale(1);
          opacity: 0.25;
        }
        25% {
          transform: translateX(7px) translateY(2px) scale(1.05);
          opacity: 0.5;
        }
        50% {
          transform: translateX(25px) translateY(-15px) scale(1.1);
          opacity: 0.8;
        }
        75% {
          transform: translateX(7px) translateY(2px) scale(1.05);
          opacity: 0.5;
        }
        100% {
          transform: translateX(-10px) translateY(20px) scale(1);
          opacity: 0.25;
        }
      }
      
      @keyframes pulseGlow {
        0% {
          filter: brightness(1) blur(0px);
          opacity: 1;
        }
        50% {
          filter: brightness(1.3) blur(0.3px);
          opacity: 1;
        }
        100% {
          filter: brightness(1) blur(0px);
          opacity: 1;
        }
      }
      
      @keyframes fadeInBackground {
        0% {
          opacity: 0;
          transform: scale(0.95);
        }
        100% {
          opacity: 1;
          transform: scale(1);
        }
      }
      
      .animate-header {
        animation: fadeInSlideDown 0.8s ease-out forwards;
      }
      
      .animate-title {
        animation: fadeInSlideUp 1s ease-out forwards;
      }
      
      .animate-subtitle {
        animation: fadeInSlideUp 1s ease-out 0.3s forwards;
      }
      
      .animate-divider {
        animation: fadeInSlideUp 1s ease-out 0.5s forwards, scaleX 0.8s ease-out 1s forwards;
      }
      
      @keyframes titleGlow {
        0%, 100% {
          text-shadow: 
            0 0 10px hsl(var(--primary) / 0.3),
            0 0 20px hsl(var(--primary) / 0.2),
            0 0 30px hsl(var(--primary) / 0.1);
        }
        50% {
          text-shadow: 
            0 0 20px hsl(var(--primary) / 0.5),
            0 0 40px hsl(var(--primary) / 0.3),
            0 0 60px hsl(var(--primary) / 0.2);
        }
      }
      
      .title-glow {
        animation: titleGlow 4s ease-in-out infinite;
        text-shadow: 
          0 0 10px hsl(var(--primary) / 0.3),
          0 0 20px hsl(var(--primary) / 0.2),
          0 0 30px hsl(var(--primary) / 0.1);
        transition: all 0.3s ease;
      }
      
      .title-glow:hover {
        text-shadow: 
          0 0 25px hsl(var(--primary) / 0.6),
          0 0 50px hsl(var(--primary) / 0.4),
          0 0 75px hsl(var(--primary) / 0.3),
          0 0 100px hsl(var(--primary) / 0.2);
        transform: scale(1.05);
      }
      
      /* Dark mode enhanced glow */
      .dark .title-glow {
        text-shadow: 
          0 0 15px hsl(var(--primary) / 0.4),
          0 0 30px hsl(var(--primary) / 0.3),
          0 0 45px hsl(var(--primary) / 0.2);
      }
      
      .dark .title-glow:hover {
        text-shadow: 
          0 0 30px hsl(var(--primary) / 0.7),
          0 0 60px hsl(var(--primary) / 0.5),
          0 0 90px hsl(var(--primary) / 0.4),
          0 0 120px hsl(var(--primary) / 0.3);
      }
      
      .animate-bg-line-1 {
        animation: fadeInBackground 3s ease-out 1s forwards, floatLine1 20s ease-in-out infinite 4s, pulseGlow 6s ease-in-out infinite 6s;
      }
      
      .animate-bg-line-2 {
        animation: fadeInBackground 3s ease-out 1.5s forwards, floatLine2 25s ease-in-out infinite 5s, pulseGlow 8s ease-in-out infinite 8s;
      }
      
      .animate-bg-line-3 {
        animation: fadeInBackground 3s ease-out 2s forwards, floatLine3 30s ease-in-out infinite 6s, pulseGlow 10s ease-in-out infinite 10s;
      }
      
      .red-line {
        background: linear-gradient(90deg, rgba(239, 68, 68, 0.6) 0%, rgba(239, 68, 68, 0.2) 50%, transparent 100%);
      }
      
      .blue-line {
        background: linear-gradient(90deg, rgba(59, 130, 246, 0.6) 0%, rgba(59, 130, 246, 0.2) 50%, transparent 100%);
      }
      
      .yellow-line {
        background: linear-gradient(90deg, rgba(234, 179, 8, 0.6) 0%, rgba(234, 179, 8, 0.2) 50%, transparent 100%);
      }
      
      .red-line-vertical {
        background: linear-gradient(180deg, rgba(239, 68, 68, 0.6) 0%, rgba(239, 68, 68, 0.2) 50%, transparent 100%);
      }
      
      .blue-line-vertical {
        background: linear-gradient(180deg, rgba(59, 130, 246, 0.6) 0%, rgba(59, 130, 246, 0.2) 50%, transparent 100%);
      }
      
      .yellow-line-vertical {
        background: linear-gradient(180deg, rgba(234, 179, 8, 0.6) 0%, rgba(234, 179, 8, 0.2) 50%, transparent 100%);
      }
    `
    document.head.appendChild(style)

    return () => {
      document.head.removeChild(style)
    }
  }, [])

  return (
    <div className="min-h-screen bg-background flex flex-col relative overflow-hidden">
      {/* Animated Background Lines */}
      <div className="absolute inset-0 pointer-events-none">
        {/* Red Lines */}
        <div className="absolute top-1/4 left-1/4 opacity-0 animate-bg-line-1">
          <div className="w-40 h-0.5 red-line rounded-full"></div>
        </div>
        
        <div className="absolute bottom-1/3 right-1/4 opacity-0 animate-bg-line-1" style={{animationDelay: '2.5s, 6s, 8s'}}>
          <div className="w-32 h-0.5 red-line rounded-full"></div>
        </div>
        
        <div className="absolute top-1/5 right-1/3 opacity-0 animate-bg-line-2" style={{animationDelay: '3s, 7s, 9s'}}>
          <div className="w-0.5 h-20 red-line-vertical rounded-full"></div>
        </div>
        
        {/* Blue Lines */}
        <div className="absolute top-1/3 right-1/5 opacity-0 animate-bg-line-2">
          <div className="w-36 h-0.5 blue-line rounded-full"></div>
        </div>
        
        <div className="absolute bottom-1/4 left-1/3 opacity-0 animate-bg-line-2" style={{animationDelay: '2.8s, 6.5s, 8.5s'}}>
          <div className="w-28 h-0.5 blue-line rounded-full"></div>
        </div>
        
        <div className="absolute top-1/2 left-1/6 opacity-0 animate-bg-line-3" style={{animationDelay: '3.2s, 7.5s, 9.5s'}}>
          <div className="w-0.5 h-24 blue-line-vertical rounded-full"></div>
        </div>
        
        {/* Yellow Lines */}
        <div className="absolute bottom-1/3 left-1/5 opacity-0 animate-bg-line-3">
          <div className="w-44 h-0.5 yellow-line rounded-full"></div>
        </div>
        
        <div className="absolute top-2/5 right-1/6 opacity-0 animate-bg-line-3" style={{animationDelay: '3.5s, 8s, 10s'}}>
          <div className="w-30 h-0.5 yellow-line rounded-full"></div>
        </div>
        
        <div className="absolute bottom-1/5 right-1/4 opacity-0 animate-bg-line-1" style={{animationDelay: '4s, 8.5s, 10.5s'}}>
          <div className="w-0.5 h-18 yellow-line-vertical rounded-full"></div>
        </div>
        
        {/* Additional Decorative Lines */}
        <div className="absolute top-1/6 left-1/2 opacity-0 animate-bg-line-2" style={{animationDelay: '4.5s, 9s, 11s'}}>
          <div className="w-16 h-0.5 red-line rounded-full"></div>
        </div>
        
        <div className="absolute bottom-1/6 left-1/4 opacity-0 animate-bg-line-1" style={{animationDelay: '5s, 9.5s, 11.5s'}}>
          <div className="w-20 h-0.5 blue-line rounded-full"></div>
        </div>
        
        <div className="absolute top-3/4 right-1/2 opacity-0 animate-bg-line-3" style={{animationDelay: '5.5s, 10s, 12s'}}>
          <div className="w-24 h-0.5 yellow-line rounded-full"></div>
        </div>
        
        {/* More Red Lines */}
        <div className="absolute top-1/8 left-3/4 opacity-0 animate-bg-line-1" style={{animationDelay: '6s, 10.5s, 12.5s'}}>
          <div className="w-26 h-0.5 red-line rounded-full"></div>
        </div>
        
        <div className="absolute bottom-1/8 left-1/6 opacity-0 animate-bg-line-2" style={{animationDelay: '6.5s, 11s, 13s'}}>
          <div className="w-0.5 h-16 red-line-vertical rounded-full"></div>
        </div>
        
        <div className="absolute top-3/5 left-2/3 opacity-0 animate-bg-line-3" style={{animationDelay: '7s, 11.5s, 13.5s'}}>
          <div className="w-22 h-0.5 red-line rounded-full"></div>
        </div>
        
        <div className="absolute bottom-2/5 right-1/8 opacity-0 animate-bg-line-1" style={{animationDelay: '7.5s, 12s, 14s'}}>
          <div className="w-0.5 h-14 red-line-vertical rounded-full"></div>
        </div>
        
        {/* More Blue Lines */}
        <div className="absolute top-1/12 right-2/5 opacity-0 animate-bg-line-2" style={{animationDelay: '8s, 12.5s, 14.5s'}}>
          <div className="w-30 h-0.5 blue-line rounded-full"></div>
        </div>
        
        <div className="absolute bottom-1/12 right-3/5 opacity-0 animate-bg-line-3" style={{animationDelay: '8.5s, 13s, 15s'}}>
          <div className="w-18 h-0.5 blue-line rounded-full"></div>
        </div>
        
        <div className="absolute top-4/5 left-1/8 opacity-0 animate-bg-line-1" style={{animationDelay: '9s, 13.5s, 15.5s'}}>
          <div className="w-0.5 h-22 blue-line-vertical rounded-full"></div>
        </div>
        
        <div className="absolute top-7/12 right-3/4 opacity-0 animate-bg-line-2" style={{animationDelay: '9.5s, 14s, 16s'}}>
          <div className="w-34 h-0.5 blue-line rounded-full"></div>
        </div>
        
        <div className="absolute bottom-7/12 left-3/5 opacity-0 animate-bg-line-3" style={{animationDelay: '10s, 14.5s, 16.5s'}}>
          <div className="w-0.5 h-26 blue-line-vertical rounded-full"></div>
        </div>
        
        {/* More Yellow Lines */}
        <div className="absolute top-5/6 right-1/6 opacity-0 animate-bg-line-1" style={{animationDelay: '10.5s, 15s, 17s'}}>
          <div className="w-28 h-0.5 yellow-line rounded-full"></div>
        </div>
        
        <div className="absolute bottom-5/6 left-2/5 opacity-0 animate-bg-line-2" style={{animationDelay: '11s, 15.5s, 17.5s'}}>
          <div className="w-32 h-0.5 yellow-line rounded-full"></div>
        </div>
        
        <div className="absolute top-1/3 left-5/6 opacity-0 animate-bg-line-3" style={{animationDelay: '11.5s, 16s, 18s'}}>
          <div className="w-0.5 h-28 yellow-line-vertical rounded-full"></div>
        </div>
        
        <div className="absolute bottom-2/3 right-2/3 opacity-0 animate-bg-line-1" style={{animationDelay: '12s, 16.5s, 18.5s'}}>
          <div className="w-38 h-0.5 yellow-line rounded-full"></div>
        </div>
        
        <div className="absolute top-11/12 left-4/5 opacity-0 animate-bg-line-2" style={{animationDelay: '12.5s, 17s, 19s'}}>
          <div className="w-0.5 h-12 yellow-line-vertical rounded-full"></div>
        </div>
        
        {/* Small Accent Lines */}
        <div className="absolute top-1/7 left-1/7 opacity-0 animate-bg-line-3" style={{animationDelay: '13s, 17.5s, 19.5s'}}>
          <div className="w-12 h-0.5 red-line rounded-full"></div>
        </div>
        
        <div className="absolute bottom-1/7 right-1/7 opacity-0 animate-bg-line-1" style={{animationDelay: '13.5s, 18s, 20s'}}>
          <div className="w-14 h-0.5 blue-line rounded-full"></div>
        </div>
        
        <div className="absolute top-6/7 left-6/7 opacity-0 animate-bg-line-2" style={{animationDelay: '14s, 18.5s, 20.5s'}}>
          <div className="w-10 h-0.5 yellow-line rounded-full"></div>
        </div>
        
        <div className="absolute bottom-6/7 right-6/7 opacity-0 animate-bg-line-3" style={{animationDelay: '14.5s, 19s, 21s'}}>
          <div className="w-0.5 h-10 red-line-vertical rounded-full"></div>
        </div>
        
        <div className="absolute top-5/12 right-5/6 opacity-0 animate-bg-line-1" style={{animationDelay: '15s, 19.5s, 21.5s'}}>
          <div className="w-0.5 h-8 blue-line-vertical rounded-full"></div>
        </div>
        
        <div className="absolute bottom-5/12 left-5/6 opacity-0 animate-bg-line-2" style={{animationDelay: '15.5s, 20s, 22s'}}>
          <div className="w-8 h-0.5 yellow-line rounded-full"></div>
        </div>
      </div>

      {/* Main Content */}
      <main className="flex-1 flex items-center justify-center px-4">
        <div className="text-center space-y-8">
          <div className="opacity-0 animate-title">
            <h1 className="text-8xl md:text-9xl lg:text-[12rem] font-bold text-foreground mb-6 tracking-tight title-glow cursor-default">
              OJ Lab
            </h1>
          </div>
          <div className="opacity-0 animate-subtitle">
            <p className="text-2xl md:text-3xl text-muted-foreground font-light tracking-wide">
              Online Judge Platform
            </p>
          </div>
          <div className="opacity-0 animate-divider">
            <div className="w-24 h-1 bg-primary mx-auto rounded-full transform scale-x-0 origin-center"></div>
          </div>
        </div>
      </main>
    </div>
  )
}
