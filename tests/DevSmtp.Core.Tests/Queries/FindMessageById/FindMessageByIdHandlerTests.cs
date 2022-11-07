using DevSmtp.Core.Models;
using DevSmtp.Core.Queries;
using DevSmtp.Core.Stores;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using Moq;

namespace DevSmtp.Core.Tests.Queries
{
    [TestClass]
    public class FindMessageByIdHandlerTests
    {
        [TestMethod]
        public async Task ExecuteAsync_WhenQueryIsValid_ItShouldFindMessage()
        {
            // Arrange
            var id = MessageId.From("id");
            var query = new FindMessageById(id);
            var to = new List<Email>();
            to.Add(Email.From("to@fake.example.com"));

            var message = new Message
            {
                Id = id,
                From = Email.From("from@fake.example.com"),
                To =to,
                Data = "message data"
            };

            // Mocks
            var dataStore = new Mock<IDataStore>(MockBehavior.Strict);
            dataStore
                .Setup(store => store.FindByIdAsync(id, default))
                .Returns(Task.FromResult<Message?>(message));

            // Act
            var handler = new FindMessageByIdHandler(dataStore.Object);
            var results = await handler.ExecuteAsync(query);

            // Assert
            Assert.IsTrue(results.Succeeded);
            Assert.IsNotNull(results.Message);
            Assert.AreEqual(id, results.Message.Id);
        }

        [TestMethod]
        public async Task ExecuteAsync_WhenQueryFails_ItShouldProduceFailureResult()
        {
            // Arrange
            var id = MessageId.From("id");
            var query = new FindMessageById(id);
            
            // Mocks
            var dataStore = new Mock<IDataStore>(MockBehavior.Strict);
            dataStore
                .Setup(store => store.FindByIdAsync(id, default))
                .Throws(new InvalidOperationException("Invalid Operation"));

            // Act
            var handler = new FindMessageByIdHandler(dataStore.Object);
            var results = await handler.ExecuteAsync(query);

            // Assert
            Assert.IsFalse(results.Succeeded);
            Assert.IsNotNull(results.Error);
            Assert.IsInstanceOfType(results.Error, typeof(FindMessageByIdException));
            Assert.IsInstanceOfType(results.Error.InnerException, typeof(InvalidOperationException));
        }
    }
}
