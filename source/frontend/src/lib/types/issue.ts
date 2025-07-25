export type Issue = {
  id: string;
  title: string;
  description: string;
  progress: string;
  date_submitted: string;
  equipment_id: string;
};

export type CreateIssueRequest = {
  title: string;
  description: string;
  equipmentId: string;
  files?: File[]; // optional media attachments
};
