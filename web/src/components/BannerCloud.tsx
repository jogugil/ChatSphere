// src/components/BannerCloud.tsx
import React, { useState } from 'react';
import { handleBannerClick, handleMouseEnter } from './banners';

interface BannerCloudProps {
  title: string;
  subtitle: string;
  imageUrl: string;
}

const BannerCloud: React.FC<BannerCloudProps> = ({ title, subtitle, imageUrl }) => {
  // Estado para mostrar u ocultar el resumen
  const [isHovered, setIsHovered] = useState(false);

  return (
    <div
      className="banner-container cloud-services"
      onClick={() => handleBannerClick('Cloud Computing')}
      onMouseEnter={() => {
        handleMouseEnter('Cloud Computing');
        setIsHovered(true); // Mostrar el resumen al pasar el mouse
      }}
      onMouseLeave={() => setIsHovered(false)} // Ocultar el resumen cuando el mouse sale
    >
      <div className="banner-content">
        <h1 className="banner-title">{title}</h1>
        <p className="banner-subtitle">{subtitle}</p>
        
        {/* Mostrar la imagen del banner */}
        <img src={imageUrl} alt="Banner de Cloud Computing" className="banner-image" />
        
        {/* Si el mouse está encima, mostrar el resumen */}
        {isHovered && (
          <div className="banner-summary">
            <h2>Temas de Cloud Computing</h2>
            <p>FaaS, Serverless, microservicios y escalabilidad automática.</p>
            <div className="tech-icons">
              <img src="../../public/images/cloudcomm.png" alt="cloud" className="icon" />
              <img src="https://es.m.wikipedia.org/wiki/Archivo:Amazon_Web_Services_Logo.svg" alt="AWS" className="icon" />
              <img src="https://logodownload.org/wp-content/uploads/2021/06/google-cloud-logo-0.png" alt="Google Cloud" className="icon" />
              <img src="https://upload.wikimedia.org/wikipedia/commons/f/fa/Microsoft_Azure.svg" alt="Azure" className="icon" />
            </div>
          </div>
 
        )}
     </div>
     </div>
  );
};

export default BannerCloud;