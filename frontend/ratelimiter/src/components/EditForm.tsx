import React from 'react';
import { Grid, TextField, Button } from '@mui/material';
import { EditFormProps } from '../types/types';

export const EditForm: React.FC<EditFormProps> = ({ newRateLimit, onRateLimitChange, onSave }) => (
  <Grid 
    container 
    spacing={2} 
    alignItems="center"
    sx={{ mt: 1 }}
  >
    <Grid item xs={8}>
      <TextField
        fullWidth
        type="number"
        size="small"
        value={newRateLimit}
        onChange={(e) => onRateLimitChange(e.target.value)}
        placeholder="New rate limit"
        sx={{
          '& .MuiOutlinedInput-root': {
            bgcolor: 'rgba(255, 255, 255, 0.03)',
            borderRadius: 1,
            border: '1px solid',
            borderColor: 'rgba(255, 255, 255, 0.1)',
            transition: 'all 0.2s ease',
            '&:hover': {
              borderColor: 'rgba(255, 255, 255, 0.2)',
            },
            '&.Mui-focused': {
              borderColor: '#5E6AD2',
              bgcolor: 'rgba(94, 106, 210, 0.04)',
            },
          },
          '& .MuiOutlinedInput-input': {
            color: 'text.primary',
            '&::placeholder': {
              color: 'rgba(255, 255, 255, 0.3)',
            },
          },
          '& .MuiOutlinedInput-notchedOutline': {
            border: 'none',
          },
        }}
      />
    </Grid>
    <Grid item xs={4}>
      <Button 
        variant="contained" 
        size="medium"
        onClick={onSave}
        fullWidth
        disableElevation
        sx={{
          height: '36px',
          bgcolor: 'rgba(255, 255, 255, 0.08)',
          color: '#fff',
          fontSize: '14px',
          fontWeight: 500,
          border: '1px solid rgba(255, 255, 255, 0.05)',
          '&:hover': {
            bgcolor: 'rgba(255, 255, 255, 0.12)',
            borderColor: 'rgba(255, 255, 255, 0.1)',
          },
          '&:active': {
            bgcolor: 'rgba(255, 255, 255, 0.15)',
          },
          '&:focus-visible': {
            outline: '2px solid #5E6AD2',
            outlineOffset: '1px',
          },
          transition: 'all 120ms ease',
          '@media (max-width: 600px)': {
            height: '40px',
          },
        }}
      >
        Save
      </Button>
    </Grid>
  </Grid>
);