// src/components/BannerProgramming.tsx
import React, { useState } from 'react';
import { handleBannerClick, handleMouseEnter } from './banners';

interface BannerProgrammingProps {
  title: string;
  subtitle: string;
  imageUrl: string;
}

const BannerProgramming: React.FC<BannerProgrammingProps> = ({ title, subtitle, imageUrl }) => {
  // Estado para mostrar u ocultar el resumen
  const [isHovered, setIsHovered] = useState(false);

  return (
    <div
      className="banner-container cloud-programming"
      onClick={() => handleBannerClick('Programación en la Nube')}
      onMouseEnter={() => {
        handleMouseEnter('Programación en la Nube');
        setIsHovered(true); // Mostrar el resumen al pasar el mouse
      }}
      onMouseLeave={() => setIsHovered(false)} // Ocultar el resumen cuando el mouse sale
    >
      <div className="banner-content">
        <h1 className="banner-title">{title}</h1>
        <p className="banner-subtitle">{subtitle}</p>
        
        {/* Mostrar la imagen del banner */}
        <img src={imageUrl} alt="Banner de Programación en la Nube" className="banner-image" />
        
        {/* Si el mouse está encima, mostrar el resumen */}
        {isHovered && (
          <div className="banner-summary">
            <h2>Temas de Programación</h2>
            <p>Patrones de diseño, concurrencia y paralelismo, escalabilidad y más.</p>
            <div className="tech-icons">
              <img src="https://miro.medium.com/v2/resize:fit:740/1*rxDdNJHiz1R38J_JEz23Zw.jpeg" alt="Go" className="icon" />
              <img src="https://upload.wikimedia.org/wikipedia/commons/d/d9/Node.js_logo.svg" alt="Node.js" className="icon" />
              <img src="https://www.typescriptlang.org/images/branding/two-longform.svg" alt="TypeScript" className="icon" />
              <img src="https://upload.wikimedia.org/wikipedia/commons/4/47/React.svg" alt="React" className="icon" />
              <img src="../../public/images/pattern.png " alt="pattern" className="icon" />
              <img src="../../public/images/concepts.png " alt="concepts" className="icon" />
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default BannerProgramming;