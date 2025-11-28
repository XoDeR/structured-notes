export interface User {
  id: string;
  username: string;
  firstname?: string;
  lastname?: string;
  role: number; // 0: User; 1: Admin;
  avatar?: string;
  password?: string;
  email: string;
  created_timestamp: number;
  updated_timestamp: number;
}

export interface PublicUser {
  id: string;
  username: string;
  avatar?: string;
  email: string;
  created_timestamp: number;
  updated_timestamp: number;
}

export interface ConnectionLog {
  id: string;
  user_id: string;
  ip_adress?: string;
  user_agent?: string;
  location?: string;
  type: string;
  timestamp: number;
}