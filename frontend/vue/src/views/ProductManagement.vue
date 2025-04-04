<template>
  <div class="product-management">
    <h1 class="page-title">Product Management</h1>
    <div class="action-bar">
      <a-button type="primary" @click="showModal" class="add-button" :style="{ backgroundColor: '#52c41a', borderColor: '#52c41a' }">
        Add Product
      </a-button>
    </div>
    <!-- ตารางสินค้า -->
    <a-table 
      :columns="columns" 
      :data-source="products" 
      row-key="id"
      class="product-table"
      :rowClassName="() => 'table-row'"
      :pagination="{ pageSize: 10 }"
      bordered
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <a-button type="primary" size="small" @click="editProduct(record)">Edit</a-button>
          <a-popconfirm
            title="Are you sure you want to delete this product?"
            ok-text="Yes"
            cancel-text="No"
            @confirm="deleteProduct(record.id)"
          >
          <a-button type="primary" danger size="small" style="margin-left: 8px">Delete</a-button>
          </a-popconfirm>
        </template>
        <template v-else-if="column.dataIndex === 'price'">
          {{ formatPrice(record.price) }}
        </template>
      </template>
    </a-table>

    <!-- Modal เพิ่ม/แก้ไขสินค้า -->
    <a-modal 
      v-model:visible="isModalVisible" 
      :title="isEditing ? 'Edit Product' : 'Add Product'" 
      @cancel="handleCancel" 
      @ok="handleSubmit"
      okText="Submit"
    >
      <a-form :model="product">
        <a-form-item label="Product Name">
          <a-input v-model:value="product.name" placeholder="Enter product name" />
        </a-form-item>
        <a-form-item label="Price">
          <!-- ใช้ a-input-number แทน a-input สำหรับตัวเลข -->
          <a-input-number v-model:value="product.price" placeholder="Enter product price" />
        </a-form-item>
        <a-form-item label="Description">
          <a-textarea v-model:value="product.description" placeholder="Enter product description" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script>
import api from '../api/product'; 

export default {
  name: 'ProductManagement',
  data() {
    return {
      products: [],
      isModalVisible: false,
      isEditing: false,
      product: { name: '', price: 0, description: '' },
      columns: [
        { title: 'Product Name', dataIndex: 'name', key: 'name' },
        { title: 'Price', dataIndex: 'price', key: 'price' },
        { title: 'Description', dataIndex: 'description', key: 'description' },
        { 
          title: 'Action', 
          key: 'action'
        }
      ],
      socket: null,
    };
  },
  methods: {
    formatPrice(price) {
      // แปลงราคาให้อยู่ในรูปแบบเงินบาท
      return '฿' + price.toLocaleString('th-TH', {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2
      });
    },
    
    async fetchProducts() {
      try {
        const fetchedProducts = await api.getProducts();
        this.products = fetchedProducts;
      } catch (error) {
        console.error("Error fetching products:", error);
      }
    },

    setupWebSocket() {
      this.socket = new WebSocket("ws://localhost:8080/ws");

      this.socket.onmessage = (event) => {
        const updatedProduct = JSON.parse(event.data);
        const index = this.products.findIndex(p => p.id === updatedProduct.id);
        if (index !== -1) {
          // ถ้าสินค้ามีอยู่ในรายการ ให้แก้ไข
          this.products[index] = updatedProduct;
        } else {
          // ถ้ายังไม่มีสินค้าในตาราง ให้เพิ่ม
          this.products.push(updatedProduct);
        }
      };

      this.socket.onopen = () => console.log("WebSocket Connected");
      this.socket.onclose = () => console.log("WebSocket Disconnected");
      this.socket.onerror = (error) => console.error("WebSocket Error:", error);
    },

    showModal() {
      this.product = { name: '', price: 0, description: '' };
      this.isEditing = false;
      this.isModalVisible = true;
    },

    editProduct(record) {
      this.product = { ...record };
      this.isEditing = true;
      this.isModalVisible = true;
    },

    handleCancel() {
      this.isModalVisible = false;
    },

    async handleSubmit() {
      if (!this.product.name || !this.product.price || !this.product.description) {
        this.$message.error("Please fill in all fields correctly");
        return;
      }

      try {
        if (this.isEditing) {
          // ถ้าเป็นการแก้ไข
          await api.updateProduct(this.product.id, this.product);
          // อัปเดตข้อมูลในอาร์เรย์ products
          const index = this.products.findIndex(p => p.id === this.product.id);
          if (index !== -1) {
            this.products[index] = { ...this.product };
          }
          this.$message.success("Product updated successfully");
        } else {
          // ถ้าเป็นการเพิ่มใหม่
          const productWithId = {
            id: Date.now().toString(),
            ...this.product
          };
          await api.addProduct(productWithId);
          this.products.push(productWithId);
          this.$message.success("Product added successfully");
        }

        this.isModalVisible = false;
      } catch (error) {
        console.error(this.isEditing ? "Failed to update product:" : "Failed to add product:", error);
        this.$message.error(this.isEditing ? "Failed to update product" : "Failed to add product");
      }
    },

    async deleteProduct(id) {
      try {
        await api.deleteProduct(id);
        this.products = this.products.filter(p => p.id !== id);
        this.$message.success(`Product deleted successfully.`);
      } catch (error) {
        console.error("Error deleting product:", error);
        this.$message.error('Failed to delete product');
      }
    }
  },

  mounted() {
    this.fetchProducts();
    this.setupWebSocket();
  },
  beforeUnmount() {
    if (this.socket) {
      this.socket.close(); 
    }
  }
};
</script>

<style scoped>
.product-management {
  padding: 20px;
}

.page-title {
  color: #000000;
  margin-bottom: 24px;
  font-size: 20px;
}

.action-bar {
  margin-bottom: 24px;
}

.add-button {
  font-weight: 500;
}

</style>