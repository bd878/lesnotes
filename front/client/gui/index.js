import React from 'react';
import { hydrateRoot } from 'react-dom/client';
import Index from './pages/Index';

hydrateRoot(document.getElementById('app'), <Index />)