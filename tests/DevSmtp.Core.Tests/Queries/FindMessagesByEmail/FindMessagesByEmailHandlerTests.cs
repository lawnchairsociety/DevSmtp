using DevSmtp.Core.Models;
using DevSmtp.Core.Queries;
using DevSmtp.Core.Stores;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using Moq;

namespace DevSmtp.Core.Tests.Queries
{
    [TestClass]
    public class FindMessagesByEmailHandlerTests
    {
        [TestMethod]
        public async Task ExecuteAsync_WhenQueryIsValid_ItShouldFindMessages()
        {
            // Arrange
            var messages = new List<Message>(100);
            var searchAddress = Email.From("from@fake.example.com");
            var to = new List<Email>();
            to.Add(Email.From("to@fake.example.com"));

            for (int i = 0; i < 100; i++)
            {
                messages.Add(new Message
                {
                    Id = MessageId.From($"{i}"),
                    To = to,
                    From = searchAddress,
                    Data = i.ToString()
                });
            }

            var query = new FindMessagesByEmail(searchAddress);

            // Mocks
            var dataStore = new Mock<IDataStore>(MockBehavior.Strict);
            dataStore
                .Setup(store => store.FindByEmailAsync(searchAddress, default))
                .Returns((Email address, CancellationToken _) =>
                {
                    var fetched = messages.AsEnumerable();
                    return Task.FromResult(fetched);
                });

            // Act
            var handler = new FindMessagesByEmailHandler(dataStore.Object);
            var results = await handler.ExecuteAsync(query);

            // Assert
            Assert.IsTrue(results.Succeeded);
            Assert.AreEqual(results.Messages.ElementAt(0).Id!.Value, "0");
            Assert.AreEqual(results.Messages.ElementAt(0).From!.Value, "from@fake.example.com");
        }

        [TestMethod]
        public async Task ExecuteAsync_WhenQueryFails_ItShouldProduceFailureResult()
        {
            // Arrange
            var searchAddress = Email.From("from@fake.example.com");
            var query = new FindMessagesByEmail(searchAddress);

            // Mocks
            var dataStore = new Mock<IDataStore>(MockBehavior.Strict);
            dataStore
            .Setup(store => store.FindByEmailAsync(searchAddress, default))
            .Throws(new InvalidOperationException("Invalid Operation"));

            // Act
            var handler = new FindMessagesByEmailHandler(dataStore.Object);
            var results = await handler.ExecuteAsync(query);

            // Assert
            Assert.IsFalse(results.Succeeded);
            Assert.IsNotNull(results.Error);
            Assert.IsInstanceOfType(results.Error, typeof(FindMessagesByEmailException));
            Assert.IsInstanceOfType(results.Error.InnerException, typeof(InvalidOperationException));
        }
    }
}
