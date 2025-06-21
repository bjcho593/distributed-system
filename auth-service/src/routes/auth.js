import express from 'express';
import jwt from 'jsonwebtoken';
import { redisClient } from '../services/redisClient.js';

const router = express.Router();

// LOGIN: genera y guarda token
router.post('/login', async (req, res) => {
  const { username, password } = req.body;

  // Validación dummy
  if (username === 'admin' && password === 'admin') {
    const token = jwt.sign({ user: username }, process.env.JWT_SECRET, { expiresIn: '1h' });

    // Guardar en Redis (clave = token, valor = username)
    await redisClient.setEx(token, 3600, username);

    return res.json({ token });
  }

  res.status(401).json({ message: 'Invalid credentials' });
});

// VALIDATE: verifica firma del token y existencia en Redis
router.get('/validate', async (req, res) => {
  try {
    const authHeader = req.headers['authorization'];
    if (!authHeader) {
      return res.status(401).json({ message: 'Missing Authorization header' });
    }

    const token = authHeader.split(' ')[1];
    const payload = jwt.verify(token, process.env.JWT_SECRET);

    const stored = await redisClient.get(token);
    if (!stored) {
      return res.status(401).json({ message: 'Token not found or expired' });
    }

    return res.json({ message: 'Token is valid', user: stored });

  } catch (err) {
    return res.status(401).json({ message: 'Invalid token', error: err.message });
  }
});

// LOGOUT: elimina token de Redis
router.post('/logout', async (req, res) => {
  try {
    const authHeader = req.headers['authorization'];
    if (!authHeader) {
      return res.status(401).json({ message: 'Missing Authorization header' });
    }

    const token = authHeader.split(' ')[1];
    jwt.verify(token, process.env.JWT_SECRET);

    const deleted = await redisClient.del(token);
    if (deleted === 0) {
      return res.status(404).json({ message: 'Token not found in Redis' });
    }

    return res.json({ message: 'Successfully logged out' });

  } catch (err) {
    return res.status(401).json({ message: 'Invalid token', error: err.message });
  }
});

export default router;
