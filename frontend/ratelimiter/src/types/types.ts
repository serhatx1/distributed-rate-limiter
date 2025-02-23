export interface PathConfig {
  id: string;
  path: string;
  requestsPerHour: number;
  type: string;
}

export interface EditFormProps {
  newRateLimit: string;
  onRateLimitChange: (value: string) => void;
  onSave: () => void;
}

export interface PathListItemProps {
  pathConfig: PathConfig;
  isEditing: boolean;
  onEditClick: () => void;
  editForm: React.ReactNode;
}
