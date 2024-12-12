import React, { createContext, useContext, useState, ReactNode } from 'react';

interface AuthContextType {
  token: string;
  nickname: string;
  roomId: string;
  roomName: string;
  setToken: (token: string) => void;
  setNickname: (nickname: string) => void;
  setRoomId: (roomId: string) => void;
  setRoomName: (roomName: string) => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [token, setToken] = useState('');
  const [nickname, setNickname] = useState('');
  const [roomId, setRoomId] = useState('');
  const [roomName, setRoomName] = useState('');

  return (
    <AuthContext.Provider value={{ token, nickname, roomId, roomName, setToken, setNickname, setRoomId, setRoomName }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
