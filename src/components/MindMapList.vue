<template>
  <div class="mindmap-list">
    <div class="list-header">
      <h2>{{ $t('mindmap.myMaps') }}</h2>
      <button @click="createNewMap" class="create-btn">
        {{ $t('mindmap.createNew') }}
      </button>
    </div>
    
    <div v-if="loading" class="loading">
      {{ $t('mindmap.loading') }}
    </div>
    
    <div v-else-if="mindMaps.length === 0" class="empty-state">
      <p>{{ $t('mindmap.noMaps') }}</p>
      <button @click="createNewMap" class="create-btn">
        {{ $t('mindmap.createFirst') }}
      </button>
    </div>
    
    <div v-else class="maps-grid">
      <div 
        v-for="map in mindMaps" 
        :key="map.id" 
        class="map-card"
        @click="openMap(map)"
      >
        <div class="map-header">
          <h3>{{ map.title }}</h3>
          <div class="map-actions">
            <button 
              @click.stop="editMap(map)"
              class="action-btn edit-btn"
              :title="$t('mindmap.edit')"
            >
              ✏️
            </button>
            <button 
              @click.stop="deleteMap(map)"
              class="action-btn delete-btn"
              :title="$t('mindmap.delete')"
            >
              🗑️
            </button>
          </div>
        </div>
        <div class="map-info">
          <span class="map-date">{{ formatDate(map.updated_at) }}</span>
          <span v-if="map.is_public" class="public-badge">
            {{ $t('mindmap.public') }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'MindMapList',
  data() {
    return {
      mindMaps: [],
      loading: true
    }
  },
  async mounted() {
    await this.loadMindMaps()
  },
  methods: {
    async loadMindMaps() {
      try {
        this.loading = true
        const response = await fetch('/api/mindmaps', {
          credentials: 'include'
        })
        
        if (response.ok) {
          //this.mindMaps = await response.json()
        } else {
          console.error('Failed to load mindmaps')
        }
      } catch (error) {
        console.error('Error loading mindmaps:', error)
      } finally {
        this.loading = false
      }
    },
    
    createNewMap() {
      this.$router.push('/edit')
    },
    
    openMap(map) {
      this.$router.push(`/edit/${map.id}`)
    },
    
    editMap(map) {
      this.$router.push(`/edit/${map.id}`)
    },
    
    async deleteMap(map) {
      if (!confirm(this.$t('mindmap.confirmDelete'))) {
        return
      }
      
      try {
        const response = await fetch(`/api/mindmaps/${map.id}`, {
          method: 'DELETE',
          credentials: 'include'
        })
        
        if (response.ok) {
          await this.loadMindMaps()
        } else {
          console.error('Failed to delete mindmap')
        }
      } catch (error) {
        console.error('Error deleting mindmap:', error)
      }
    },
    
    formatDate(dateString) {
      const date = new Date(dateString)
      return date.toLocaleDateString()
    }
  }
}
</script>

<style scoped>
.mindmap-list {
  padding: 2rem;
  max-width: 1200px;
  margin: 0 auto;
}

.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.list-header h2 {
  margin: 0;
  color: #333;
}

.create-btn {
  background: #667eea;
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
  transition: background-color 0.2s;
}

.create-btn:hover {
  background: #5a6fd8;
}

.loading {
  text-align: center;
  padding: 2rem;
  color: #666;
}

.empty-state {
  text-align: center;
  padding: 4rem 2rem;
  color: #666;
}

.empty-state p {
  margin-bottom: 1rem;
  font-size: 1.1rem;
}

.maps-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1.5rem;
}

.map-card {
  background: white;
  border: 1px solid #e1e5e9;
  border-radius: 8px;
  padding: 1.5rem;
  cursor: pointer;
  transition: all 0.2s;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.map-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
  border-color: #667eea;
}

.map-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1rem;
}

.map-header h3 {
  margin: 0;
  color: #333;
  font-size: 1.1rem;
  flex: 1;
}

.map-actions {
  display: flex;
  gap: 0.5rem;
}

.action-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 0.25rem;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.action-btn:hover {
  background: #f0f0f0;
}

.delete-btn:hover {
  background: #fee;
}

.map-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.9rem;
  color: #666;
}

.public-badge {
  background: #e8f5e8;
  color: #2d5a2d;
  padding: 0.25rem 0.5rem;
  border-radius: 12px;
  font-size: 0.8rem;
}
</style> 