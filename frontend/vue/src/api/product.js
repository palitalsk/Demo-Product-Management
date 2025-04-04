import axios from 'axios';

const API_URL = 'http://localhost:8080/products';

export default {
  // ฟังก์ชันดึงข้อมูลสินค้า
  async getProducts() {
    try {
      const res = await axios.get(API_URL);
      return res.data;  
    } catch (error) {
      console.error("Error fetching products:", error);
      throw new Error('Failed to fetch products');
    }
  },

  // ฟังก์ชันเพิ่มสินค้า
  async addProduct(product) {
    // ตรวจสอบว่ามีข้อมูลครบถ้วนก่อนส่งไปยัง backend
    if (!product.name || !product.price || !product.description) {
      throw new Error('Please fill in all fields');
    }

    if (product.price <= 0) {
      throw new Error('Price must be greater than 0');
    }

    try {
      // ส่งข้อมูลไปยัง backend
      const res = await axios.post(API_URL, product);
      console.log("Product added successfully:", res.data);
      return res.data; 
    } catch (error) {
      console.error("Error adding product:", error);
      throw new Error('Failed to add product');
    }
  },

  // ฟังก์ชันลบสินค้า
  async deleteProduct(id) {
    try {
      const res = await axios.delete(`${API_URL}/${id}`);
      console.log("Product deleted successfully:", res.data);
      return res.data; 
    } catch (error) {
      console.error("Error deleting product:", error);
      throw new Error('Failed to delete product');
    }
  },

  // ฟังก์ชันอัปเดตข้อมูลสินค้า
  async updateProduct(id, product) {
    if (!product.name || !product.price || !product.description) {
      throw new Error('Please fill in all fields');
    }

    if (product.price <= 0) {
      throw new Error('Price must be greater than 0');
    }

    try {
      const res = await axios.put(`${API_URL}/${id}`, product);
      console.log("Product updated successfully:", res.data);
      return res.data;  
    } catch (error) {
      console.error("Error updating product:", error);
      throw new Error('Failed to update product');
    }
  }

};
