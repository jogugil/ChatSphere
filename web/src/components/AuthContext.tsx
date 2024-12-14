import React, { createContext, useContext, useState, ReactNode, useMemo } from 'react';

interface AuthContextType {
  token: string;
  nickName: string;
  roomId: string;
  roomName: string;
  setToken: (token: string) => void;
  setNickName: (nickName: string) => void;
  setRoomId: (roomId: string) => void;
  setRoomName: (roomName: string) => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [token, setToken] = useState('');
  const [nickName, setNickName] = useState('');
  const [roomId, setRoomId] = useState('');
  const [roomName, setRoomName] = useState('');

  // Memorizamos el valor del contexto para evitar re-renderizados innecesarios
  const value = useMemo(() => ({
    token, 
    nickName, 
    roomId, 
    roomName, 
    setToken, 
    setNickName, 
    setRoomId, 
    setRoomName
  }), [token, nickName, roomId, roomName]); // Dependencias para re-memorizar solo cuando cambia el estado

  // Agregar console.log aquí para ver los valores en la consola
  console.log('AuthContext State:', { token, nickName, roomId, roomName });

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  console.log('useAuth State context:', { context });

  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }

  return context;
};
