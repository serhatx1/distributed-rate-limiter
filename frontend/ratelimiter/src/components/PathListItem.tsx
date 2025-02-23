import React from 'react';
import { ListItem, ListItemText, IconButton, Chip, Box, Typography } from '@mui/material';
import EditIcon from '@mui/icons-material/Edit';
import { PathListItemProps } from '../types/types';

export const PathListItem: React.FC<PathListItemProps> = ({ 
  pathConfig, 
  isEditing, 
  onEditClick, 
  editForm 
}) => (
  <ListItem
    sx={{
      borderBottom: '1px solid',
      borderColor: 'divider',
      '&:last-child': {
        borderBottom: 'none',
      },
      py: 2.5,
      px: 3.5,
      '&:hover': {
        bgcolor: 'rgba(255, 255, 255, 0.02)',
      },
      transition: 'background-color 120ms ease',
      display: 'flex',
      justifyContent: 'space-between',
      alignItems: 'flex-start',
    }}
  >
    <Box sx={{ flex: 1, mr: 2 }}>
      <Typography 
        variant="subtitle1" 
        sx={{ 
          fontWeight: 500,
          color: 'text.primary',
          fontSize: '14px',
          letterSpacing: '-0.01em',
        }}
      >
        {pathConfig.path}
      </Typography>
      <Box sx={{ mt: 0.5 }}>
        {isEditing ? (
          editForm
        ) : (
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <Chip
              label={`${pathConfig.requestsPerHour} requests/hour`}
              size="small"
              sx={{
                height: '24px',
                bgcolor: 'rgba(255, 255, 255, 0.04)',
                color: 'rgba(255, 255, 255, 0.8)',
                border: '1px solid rgba(255, 255, 255, 0.1)',
                '& .MuiChip-label': {
                  px: 1,
                  fontSize: '12px',
                  fontWeight: 500,
                },
              }}
            />
            {pathConfig?.type.toLowerCase() === "default" && (
              <Chip
                label="Default"
                size="small"
                sx={{
                  height: '24px',
                  bgcolor: 'rgba(94, 106, 210, 0.1)',
                  color: '#5E6AD2',
                  border: '1px solid rgba(94, 106, 210, 0.2)',
                  '& .MuiChip-label': {
                    px: 1,
                    fontSize: '12px',
                    fontWeight: 500,
                  },
                }}
              />
            )}
          </Box>
        )}
      </Box>
    </Box>
    <IconButton 
      onClick={onEditClick}
      sx={{ 
        color: 'rgba(255, 255, 255, 0.4)',
        padding: '6px',
        '&:hover': {
          bgcolor: 'rgba(255, 255, 255, 0.04)',
          color: 'rgba(255, 255, 255, 0.9)',
        },
      }}
    >
      <EditIcon sx={{ fontSize: 18 }} />
    </IconButton>
  </ListItem>
);