export type User = {
  id: string;
  username: string;
  email: string;
  isActive: boolean;
};

export type UserCreatePayload = {
  username: string;
  email: string;
  password: string;
};
