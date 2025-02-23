import { useState, useEffect } from 'react';
import { Container, Paper, Typography, List } from '@mui/material';
import { PathListItem } from '../components/PathListItem';
import { EditForm } from '../components/EditForm';
import { PathConfig } from '../types/types';
import { ThemeProvider } from '@mui/material/styles';
import { theme } from '../theme/theme';

export const Landing = () => {
  const [paths, setPaths] = useState<PathConfig[]>([]);
  const [editingPath, setEditingPath] = useState<string | null>(null);
  const [newRateLimit, setNewRateLimit] = useState('');

  useEffect(() => {
    fetchRoutes();
  }, []);

  const fetchRoutes = async () => {
    try {
      console.log("formattedPaths")

      const response = await fetch('http://localhost:3000/listroutes');
      const data = await response.json();
      console.log("data",data)
      if (data.routes) {
        const formattedPaths = Object.entries(data.routes).map(([path, config]: [string, any]) => ({
          id: path||"",
          path:config.path||"",
          requestsPerHour: config.rate_limit || 0,
          type:config.source||"",

        }));
        setPaths(formattedPaths);
        console.log(formattedPaths)
      }
    } catch (error) {
      console.error('Error fetching routes:', error);
    }
  };

  const handleUpdateRateLimit = async (path: string,type:string) => {
    let url = "http://localhost:3000/ratelimit/changelimit"
    if(type?.toLowerCase()==="default"){
       url = "http://localhost:3000/ratelimit/setlimit"

    }
    try {
      await fetch(url, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          "EndPoint":path,
          "Ratelimit": parseInt(newRateLimit),
        }),
      });
      await fetchRoutes(); 
      setEditingPath(null);
      setNewRateLimit('');
    } catch (error) {
      console.error('Error updating rate limit:', error);
    }
  };

  const handleEditClick = (path: string, currentRate: number) => {
    if (editingPath === path) {
      // If clicking the same path's edit button, close it
      setEditingPath(null);
      setNewRateLimit('');
    } else {
      // If clicking a different path's edit button, open it
      setEditingPath(path);
      setNewRateLimit(currentRate.toString());
    }
  };

  return (
    <ThemeProvider theme={theme}>
      <Container 
        maxWidth="md" 
        sx={{ 
          mt: 4, 
          mb: 8,
          bgcolor: '#2C2C2E',
          padding:"50px",
          borderRadius:"30px"

        }}
      >
        <Typography 
          variant="h4" 
          gutterBottom 
          sx={{ 
            mb: 4,
            fontWeight: 600,
            textAlign: 'center',
            color: 'text.primary',
            
            letterSpacing: '-0.5px',
          }}
        >
          Rate Limit Configuration
        </Typography>
        <Paper 
          elevation={0}
          sx={{ 
            borderRadius: 3,
            overflow: 'hidden',
            backgroundColor: '#1C1C1E',
            border: '1px solid',
            borderColor: 'divider',
          }}
        >
          <List sx={{ p: 0 }}>
            {paths.map((pathConfig) => (
              <PathListItem
                key={pathConfig.id}
                pathConfig={pathConfig}
                isEditing={editingPath === pathConfig.path}
                onEditClick={() => handleEditClick(pathConfig.path, pathConfig.requestsPerHour)}
                editForm={
                  <EditForm
                    newRateLimit={newRateLimit}
                    onRateLimitChange={setNewRateLimit}
                    onSave={() => handleUpdateRateLimit(pathConfig.path, pathConfig.type)}
                  />
                }
              />
            ))}
          </List>
        </Paper>
      </Container>
    </ThemeProvider>
  );
};
