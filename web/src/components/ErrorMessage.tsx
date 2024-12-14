import React, { useState } from 'react';
import '../styles/errorMessage.css';

interface WErrorMessageProps {
  message: string;
  showError: boolean;
  closeErrorMessage: () => void;
  minimizeErrorMessage: () => void;
  restoreErrorMessage: () => void;
  minimized: boolean;
  iconType: "error" | "info"; // Mantén el tipo de icono como string literal
}

// ErrorMessage.tsx
export const WErrorMessage = ({
  message,
  showError,
  closeErrorMessage,
  minimizeErrorMessage,
  restoreErrorMessage,
  minimized,
  iconType,  // Asegúrate de que iconType está definido
}: WErrorMessageProps) => {
  return (
    showError && (
      <div className={`w-message-window ${iconType === "error" ? "error" : "info"}`}>
        <div className="w-header">
          {/* Icono dinámico basado en el tipo */}
          <span className={`w-icon ${iconType}`}>
            {iconType === "error" ? (
              <i className="fas fa-times-circle"></i> // Icono de error
            ) : (
              <i className="fas fa-exclamation-triangle"></i> // Icono de información
            )}
          </span>
          <span className="w-title">
            {iconType === "error" ? "Error" : "Información"}
          </span>
          <div className="w-controls">
            <button
              onClick={minimized ? restoreErrorMessage : minimizeErrorMessage}
              title={minimized ? "Maximizar" : "Minimizar"}
            >
              {minimized ? "⬜" : "_"}
            </button>
            <button onClick={closeErrorMessage} title="Cerrar">
              X
            </button>
          </div>
        </div>
        {!minimized && (
          <div className="w-content">
            <p>{message}</p>
            <button className="w-ok-button" onClick={closeErrorMessage}>
              Ok
            </button>
          </div>
        )}
      </div>
    )
  );
};