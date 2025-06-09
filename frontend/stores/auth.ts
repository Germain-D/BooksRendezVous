import { defineStore } from 'pinia';

interface AuthState {
  isAuthenticated: boolean;
  error: string | null;
  pseudo: string | null;
  uuid: string | null;
  ispublic: boolean;
  sharelink: string | null;
  email: string | null;
}

interface Credentials {
  email: string;
  password: string;
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    isAuthenticated: false,
    error: null,
    pseudo: null,
    email: null,
    uuid: null,
    ispublic: false,
    sharelink: null,
  }),

  actions: {
    async login(credentials: Credentials) {
      try {
        const config = useRuntimeConfig();
        const response = await fetch(config.public.BACKEND_URL + '/api/login', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include', // Important for cookies
          body: JSON.stringify(credentials),
        });

        const data = await response.json();
        
        if (!response.ok) {
          this.error = data.error || 'Login failed';
          this.isAuthenticated = false;
          throw new Error(this.error || 'Unknown error');
        }

        this.isAuthenticated = true;
        this.pseudo = data.pseudo;
        this.uuid = data.uuid;
        this.email = data.email;

        // write the token in local storage
        this.setToken(data.token);


        this.error = null;
        return true;
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Unknown error';
        this.isAuthenticated = false;
        return false;
      }
    },

    async register(credentials: { username: string; email: string; password: string }) {
      try {
        const config = useRuntimeConfig();
        const response = await fetch(`${config.public.BACKEND_URL}/api/register`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(credentials),
        })

        if (!response.ok) {
          throw new Error('Registration failed')
        }
        
        return true
      } catch (error) {
        console.error('Registration error:', error)
        return false
      }
    },

   

    setToken(token: string | null) {
      this.isAuthenticated = !!token;
      if (token) {
        localStorage.setItem('jwt', token);
      } else {
        localStorage.removeItem('jwt');
      }
    },



    clearToken() {
      this.setToken(null);
    },

    async logout() {
      try {
        const config = useRuntimeConfig();
        await fetch(config.public.BACKEND_URL+'/logout', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + localStorage.getItem('jwt')
          },
          credentials: 'include',
        });
      } finally {
        this.setToken(null);
        this.isAuthenticated = false;
        this.pseudo = null;
        this.uuid = null;
      }
    },

    async checkAuth() {
      try {
        const config = useRuntimeConfig();
        const response = await fetch(config.public.BACKEND_URL+'/api/user', {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + localStorage.getItem('jwt')
          },
          credentials: 'include',
        });
        console.log(response);
        if (response.ok) {
          const data = await response.json();
          this.isAuthenticated = true;
          return true;
        }
      } catch (error) {
        this.setToken(null);
      }
      return false;
    },
    async changePublicVisibility() {
      try {
          const config = useRuntimeConfig();
          
          await fetch(config.public.BACKEND_URL + '/api/changepublicvisibility', {
              method: 'POST',
              headers: {
                  'Content-Type': 'application/json',
                  'Authorization': 'Bearer ' + localStorage.getItem('jwt')
              },
              credentials: 'include',
          });
          this.ispublic = !this.ispublic;
      } catch (error) {
          console.error('Error while fetching books', error);
      }
    },

    async getPublicVisibility() {
        try {
            const config = useRuntimeConfig();
            const response = await fetch(config.public.BACKEND_URL + '/api/getpublicvisibility', {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': 'Bearer ' + localStorage.getItem('jwt')
                },
                credentials: 'include',
            }).then(res => res.json());
            this.sharelink = response.sharelink;
            return response.public;
        } catch (error) {
            console.error('Error while fetching books', error);
        }
    },

    async passwordChange(oldpassword: string, newpassword: string): Promise<boolean> {
      try {
        const credentials = {
            oldpassword: oldpassword,
            newpassword: newpassword
        }
          const config = useRuntimeConfig();
          const response = await fetch(config.public.BACKEND_URL + '/api/passwordchange', {
              method: 'POST',
              headers: {
                  'Content-Type': 'application/json',
                  'Authorization': 'Bearer ' + localStorage.getItem('jwt')
              },
              credentials: 'include',
              body: JSON.stringify(credentials),
          });

          if (!response.ok) {
              const error = await response.json();
              throw new Error(error.message || 'Failed to change password');
          }

          return true;
      } catch (error) {
          console.error('Error while changing password:', error);
          return false;
  
          
      }
  },
  async sendResetPasswordEmail(email: string): Promise<boolean> {
    try {
        const credentials = {
            email: email
        }
        const config = useRuntimeConfig();
        const response = await fetch(config.public.BACKEND_URL + '/api/forgetpassword', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
            body: JSON.stringify(credentials),
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.message || 'Failed to change password');
        }

        return true;
    } catch (error) {
        console.error('Error while changing password:', error);
        return false;   
    }
  },

  async resetPassword(token: string, password: string): Promise<boolean> {
    try {
        const credentials = {
            token: token,
            password: password
        }
        const config = useRuntimeConfig();
        const response = await fetch(config.public.BACKEND_URL + '/api/reset-password', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
            body: JSON.stringify(credentials),
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.message || 'Failed to change password');
        }

        return true;
    } catch (error) {
        console.error('Error while changing password:', error);
        return false;
    }
  },

  verifyResetToken(token: string): Promise<boolean> {
    return new Promise((resolve, reject) => {
        const config = useRuntimeConfig();
        fetch(config.public.BACKEND_URL + '/api/verify-reset-token/', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
            body: JSON.stringify({
                token: token
            })

        }).then(response => {
            if (response.ok) {
                resolve(true);
            } else {
                resolve(false);
            }
        }).catch(error => {
            console.error('Error while verifying token:', error);
            reject(error);
        }
        );
    });
}


},

  

  getters: {

    getError: (state) => state.error,
    isLoggedIn: (state) => state.isAuthenticated,
    getPseudo: (state) => state.pseudo,
    getUuid: (state) => state.uuid,
    getIsPublic: (state) => state.ispublic,
    getShareLink: (state) => state.sharelink,
    getEmail: (state) => state.email,
  },
  persist: true,
});